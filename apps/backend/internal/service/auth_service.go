package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	maxFailedAttempts = 5
	lockDuration      = 15 * time.Minute
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Login(ctx context.Context, input domain.LoginInput, ipAddress, userAgent string) (*domain.User, *domain.AuthTokens, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthTokens, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, input domain.ChangePasswordInput) error
}

// authService implements AuthService
type authService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
	logger     zerolog.Logger
}

// NewAuthService creates a new authService
func NewAuthService(userRepo repository.UserRepository, jwtManager *auth.JWTManager, logger zerolog.Logger) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Login authenticates a user and returns tokens
func (s *authService) Login(ctx context.Context, input domain.LoginInput, ipAddress, userAgent string) (*domain.User, *domain.AuthTokens, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			// Don't reveal whether email exists
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, fmt.Errorf("authService.Login find user: %w", err)
	}

	// Check if account is active
	if !user.IsActive() {
		s.logger.Warn().
			Str("email", input.Email).
			Str("status", string(user.Status)).
			Msg("login attempt on inactive account")
		return nil, nil, domain.ErrAccountInactive
	}

	// Check if account is locked
	if user.IsLocked() {
		s.logger.Warn().
			Str("email", input.Email).
			Time("locked_until", *user.LockedUntil).
			Msg("login attempt on locked account")
		return nil, nil, domain.ErrAccountLocked
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		// Increment failed attempts
		if incrementErr := s.userRepo.IncrementFailedAttempts(ctx, user.ID); incrementErr != nil {
			s.logger.Error().Err(incrementErr).Msg("failed to increment failed attempts")
		}

		// Lock account if too many failures
		if user.FailedAttempts+1 >= maxFailedAttempts {
			lockUntil := time.Now().Add(lockDuration)
			if lockErr := s.userRepo.LockAccount(ctx, user.ID, lockUntil); lockErr != nil {
				s.logger.Error().Err(lockErr).Msg("failed to lock account")
			}
			s.logger.Warn().
				Str("email", input.Email).
				Int("attempts", user.FailedAttempts+1).
				Msg("account locked due to too many failed attempts")
		}

		return nil, nil, domain.ErrInvalidCredentials
	}

	// Reset failed attempts on successful login
	if err := s.userRepo.ResetFailedAttempts(ctx, user.ID); err != nil {
		s.logger.Error().Err(err).Msg("failed to reset failed attempts")
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID, ipAddress); err != nil {
		s.logger.Error().Err(err).Msg("failed to update last login")
	}

	// Generate tokens
	tokens, err := s.generateTokens(ctx, user, ipAddress, userAgent)
	if err != nil {
		return nil, nil, fmt.Errorf("authService.Login generate tokens: %w", err)
	}

	s.logger.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Str("ip", ipAddress).
		Msg("user logged in successfully")

	return user, tokens, nil
}

// Logout invalidates the refresh token
func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := auth.HashToken(refreshToken)
	if err := s.userRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("authService.Logout: %w", err)
	}
	return nil
}

// RefreshTokens generates new access and refresh tokens
func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (*domain.AuthTokens, error) {
	// Validate the refresh token JWT
	claims, err := s.jwtManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if token exists and is not revoked in DB
	tokenHash := auth.HashToken(refreshToken)
	storedToken, err := s.userRepo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("authService.RefreshTokens find token: %w", err)
	}

	if !storedToken.IsValid() {
		return nil, domain.ErrTokenRevoked
	}

	// Get the user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens find user: %w", err)
	}

	if !user.IsActive() {
		return nil, domain.ErrAccountInactive
	}

	// Revoke old refresh token (token rotation)
	if err := s.userRepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		s.logger.Error().Err(err).Msg("failed to revoke old refresh token")
	}

	// Generate new tokens
	ipAddress := ""
	userAgent := ""
	if storedToken.IPAddress != nil {
		ipAddress = *storedToken.IPAddress
	}
	if storedToken.UserAgent != nil {
		userAgent = *storedToken.UserAgent
	}

	tokens, err := s.generateTokens(ctx, user, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("authService.RefreshTokens generate tokens: %w", err)
	}

	return tokens, nil
}

// ChangePassword changes a user's password
func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, input domain.ChangePasswordInput) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("authService.ChangePassword find user: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return domain.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("authService.ChangePassword hash: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("authService.ChangePassword update: %w", err)
	}

	// Revoke all refresh tokens (force re-login on all devices)
	if err := s.userRepo.RevokeAllUserRefreshTokens(ctx, userID); err != nil {
		s.logger.Error().Err(err).Msg("failed to revoke all refresh tokens after password change")
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Msg("password changed successfully")

	return nil
}

// generateTokens creates access and refresh tokens and stores the refresh token
func (s *authService) generateTokens(ctx context.Context, user *domain.User, ipAddress, userAgent string) (*domain.AuthTokens, error) {
	// Generate access token
	accessToken, expiresAt, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generateTokens access: %w", err)
	}

	// Generate refresh token
	refreshTokenStr, refreshExpiresAt, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("generateTokens refresh: %w", err)
	}

	// Store refresh token hash in DB
	tokenHash := auth.HashToken(refreshTokenStr)
	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: refreshExpiresAt,
	}
	if ipAddress != "" {
		refreshToken.IPAddress = &ipAddress
	}
	if userAgent != "" {
		refreshToken.UserAgent = &userAgent
	}

	if err := s.userRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("generateTokens store refresh: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}
