package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// ─── Mock UserRepository ──────────────────────────────────────────────────────

type mockUserRepository struct {
	users         map[string]*domain.User
	refreshTokens map[string]*domain.RefreshToken
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:         make(map[string]*domain.User),
		refreshTokens: make(map[string]*domain.RefreshToken),
	}
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepository) FindAll(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int, error) {
	var users []*domain.User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, len(users), nil
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for email, u := range m.users {
		if u.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockUserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error {
	return nil
}

func (m *mockUserRepository) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	for _, u := range m.users {
		if u.ID == id {
			u.FailedAttempts++
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) ResetFailedAttempts(ctx context.Context, id uuid.UUID) error {
	for _, u := range m.users {
		if u.ID == id {
			u.FailedAttempts = 0
			u.LockedUntil = nil
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) LockAccount(ctx context.Context, id uuid.UUID, until time.Time) error {
	for _, u := range m.users {
		if u.ID == id {
			u.LockedUntil = &until
			return nil
		}
	}
	return nil
}

func (m *mockUserRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	m.refreshTokens[token.TokenHash] = token
	return nil
}

func (m *mockUserRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	if t, ok := m.refreshTokens[tokenHash]; ok {
		return t, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if t, ok := m.refreshTokens[tokenHash]; ok {
		t.IsRevoked = true
		return nil
	}
	return nil
}

func (m *mockUserRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	for _, t := range m.refreshTokens {
		if t.UserID == userID {
			t.IsRevoked = true
		}
	}
	return nil
}

func (m *mockUserRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return nil
}

// ─── Test helpers ─────────────────────────────────────────────────────────────

func createTestUser(email, password string, role domain.UserRole) *domain.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return &domain.User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  string(hash),
		FullName:      "Test User",
		Role:          role,
		Status:        domain.StatusActive,
		EmailVerified: true,
	}
}

func createTestAuthService(repo *mockUserRepository) service.AuthService {
	jwtManager := auth.NewJWTManager(
		"test-access-secret-key-minimum-32-chars",
		"test-refresh-secret-key-minimum-32-chars",
		15*time.Minute,
		7*24*time.Hour,
		"test-issuer",
	)
	logger := zerolog.Nop()
	return service.NewAuthService(repo, jwtManager, logger)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestAuthService_Login_Success(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	input := domain.LoginInput{
		Email:    "admin@test.com",
		Password: "password123",
	}

	returnedUser, tokens, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if returnedUser == nil {
		t.Fatal("expected user, got nil")
	}
	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}
	if tokens.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if tokens.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if returnedUser.Email != "admin@test.com" {
		t.Errorf("expected email admin@test.com, got %s", returnedUser.Email)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	input := domain.LoginInput{
		Email:    "admin@test.com",
		Password: "wrongpassword",
	}

	_, _, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := newMockUserRepository()
	svc := createTestAuthService(repo)

	input := domain.LoginInput{
		Email:    "nonexistent@test.com",
		Password: "password123",
	}

	_, _, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	user.Status = domain.StatusInactive
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	input := domain.LoginInput{
		Email:    "admin@test.com",
		Password: "password123",
	}

	_, _, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")

	if !errors.Is(err, domain.ErrAccountInactive) {
		t.Errorf("expected ErrAccountInactive, got: %v", err)
	}
}

func TestAuthService_Login_LockedAccount(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	lockUntil := time.Now().Add(15 * time.Minute)
	user.LockedUntil = &lockUntil
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	input := domain.LoginInput{
		Email:    "admin@test.com",
		Password: "password123",
	}

	_, _, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")

	if !errors.Is(err, domain.ErrAccountLocked) {
		t.Errorf("expected ErrAccountLocked, got: %v", err)
	}
}

func TestAuthService_Logout_Success(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	// Login first to get tokens
	input := domain.LoginInput{Email: "admin@test.com", Password: "password123"}
	_, tokens, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Logout
	err = svc.Logout(context.Background(), tokens.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error on logout, got: %v", err)
	}

	// Verify token is revoked
	tokenHash := auth.HashToken(tokens.RefreshToken)
	storedToken, _ := repo.FindRefreshToken(context.Background(), tokenHash)
	if storedToken != nil && !storedToken.IsRevoked {
		t.Error("expected refresh token to be revoked after logout")
	}
}

func TestAuthService_RefreshTokens_Success(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	// Login first
	input := domain.LoginInput{Email: "admin@test.com", Password: "password123"}
	_, tokens, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Refresh tokens
	newTokens, err := svc.RefreshTokens(context.Background(), tokens.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error on refresh, got: %v", err)
	}

	if newTokens.AccessToken == "" {
		t.Error("expected non-empty new access token")
	}
	if newTokens.RefreshToken == "" {
		t.Error("expected non-empty new refresh token")
	}
	// New tokens should be different from old ones
	if newTokens.AccessToken == tokens.AccessToken {
		t.Error("expected new access token to be different from old one")
	}
}

func TestAuthService_RefreshTokens_RevokedToken(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "password123", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	// Login first
	input := domain.LoginInput{Email: "admin@test.com", Password: "password123"}
	_, tokens, err := svc.Login(context.Background(), input, "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	// Logout (revoke token)
	svc.Logout(context.Background(), tokens.RefreshToken)

	// Try to refresh with revoked token
	_, err = svc.RefreshTokens(context.Background(), tokens.RefreshToken)
	if err == nil {
		t.Error("expected error when refreshing with revoked token")
	}
}

func TestAuthService_ChangePassword_Success(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "oldpassword", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	input := domain.ChangePasswordInput{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword123",
	}

	err := svc.ChangePassword(context.Background(), user.ID, input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify new password works
	loginInput := domain.LoginInput{Email: "admin@test.com", Password: "newpassword123"}
	_, _, err = svc.Login(context.Background(), loginInput, "127.0.0.1", "test-agent")
	if err != nil {
		t.Errorf("expected login with new password to succeed, got: %v", err)
	}
}

func TestAuthService_ChangePassword_WrongCurrentPassword(t *testing.T) {
	repo := newMockUserRepository()
	user := createTestUser("admin@test.com", "oldpassword", domain.RoleAdmin)
	repo.users[user.Email] = user

	svc := createTestAuthService(repo)

	input := domain.ChangePasswordInput{
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword123",
	}

	err := svc.ChangePassword(context.Background(), user.ID, input)
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}
