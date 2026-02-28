package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
)

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
}

// AccessClaims holds the JWT claims for access tokens
type AccessClaims struct {
	UserID uuid.UUID       `json:"user_id"`
	Email  string          `json:"email"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// RefreshClaims holds the JWT claims for refresh tokens
type RefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWTManager
func NewJWTManager(
	accessSecret, refreshSecret string,
	accessExpiry, refreshExpiry time.Duration,
	issuer string,
) *JWTManager {
	return &JWTManager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
	}
}

// GenerateAccessToken generates a new JWT access token for the given user
func (m *JWTManager) GenerateAccessToken(user *domain.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.accessExpiry)

	claims := &AccessClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.accessSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("GenerateAccessToken: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new JWT refresh token for the given user
func (m *JWTManager) GenerateRefreshToken(user *domain.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(m.refreshExpiry)

	claims := &RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.refreshSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("GenerateRefreshToken: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateAccessToken validates a JWT access token and returns the claims
func (m *JWTManager) ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.accessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrInvalidToken
		}
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates a JWT refresh token and returns the claims
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.refreshSecret), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return claims, nil
}

// HashToken creates a SHA-256 hash of a token string for secure storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
