package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
)

func createTestJWTManager() *auth.JWTManager {
	return auth.NewJWTManager(
		"test-access-secret-key-minimum-32-chars!!",
		"test-refresh-secret-key-minimum-32-chars!",
		15*time.Minute,
		7*24*time.Hour,
		"test-issuer",
	)
}

func createTestUser() *domain.User {
	return &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  domain.RoleAdmin,
	}
}

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	manager := createTestJWTManager()
	user := createTestUser()

	token, expiresAt, err := manager.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if expiresAt.Before(time.Now()) {
		t.Error("expected expiry to be in the future")
	}
	if expiresAt.After(time.Now().Add(20 * time.Minute)) {
		t.Error("expected expiry to be within 20 minutes")
	}
}

func TestJWTManager_ValidateAccessToken_Valid(t *testing.T) {
	manager := createTestJWTManager()
	user := createTestUser()

	token, _, err := manager.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	claims, err := manager.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if claims.UserID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, claims.UserID)
	}
	if claims.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, claims.Email)
	}
	if claims.Role != user.Role {
		t.Errorf("expected role %s, got %s", user.Role, claims.Role)
	}
}

func TestJWTManager_ValidateAccessToken_InvalidToken(t *testing.T) {
	manager := createTestJWTManager()

	_, err := manager.ValidateAccessToken("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestJWTManager_ValidateAccessToken_WrongSecret(t *testing.T) {
	manager1 := createTestJWTManager()
	manager2 := auth.NewJWTManager(
		"different-access-secret-key-minimum-32-chars",
		"different-refresh-secret-key-minimum-32-chars",
		15*time.Minute,
		7*24*time.Hour,
		"test-issuer",
	)

	user := createTestUser()
	token, _, _ := manager1.GenerateAccessToken(user)

	_, err := manager2.ValidateAccessToken(token)
	if err == nil {
		t.Error("expected error when validating with wrong secret")
	}
}

func TestJWTManager_RefreshToken_CannotValidateAsAccess(t *testing.T) {
	manager := createTestJWTManager()
	user := createTestUser()

	// Generate refresh token
	refreshToken, _, err := manager.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("generate refresh token failed: %v", err)
	}

	// Try to validate refresh token as access token (should fail)
	_, err = manager.ValidateAccessToken(refreshToken)
	if err == nil {
		t.Error("expected error when validating refresh token as access token")
	}
}

func TestJWTManager_AccessToken_CannotValidateAsRefresh(t *testing.T) {
	manager := createTestJWTManager()
	user := createTestUser()

	// Generate access token
	accessToken, _, err := manager.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("generate access token failed: %v", err)
	}

	// Try to validate access token as refresh token (should fail)
	_, err = manager.ValidateRefreshToken(accessToken)
	if err == nil {
		t.Error("expected error when validating access token as refresh token")
	}
}

func TestHashToken_Deterministic(t *testing.T) {
	token := "test-token-value"
	hash1 := auth.HashToken(token)
	hash2 := auth.HashToken(token)

	if hash1 != hash2 {
		t.Error("expected same hash for same token")
	}
}

func TestHashToken_DifferentTokens(t *testing.T) {
	hash1 := auth.HashToken("token1")
	hash2 := auth.HashToken("token2")

	if hash1 == hash2 {
		t.Error("expected different hashes for different tokens")
	}
}

func TestHashToken_NotEmpty(t *testing.T) {
	hash := auth.HashToken("some-token")
	if hash == "" {
		t.Error("expected non-empty hash")
	}
}
