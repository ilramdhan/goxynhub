package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/config"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/middleware"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/response"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/service"
)

const refreshTokenCookieName = "refresh_token"

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthService
	cfg         *config.Config
	logger      zerolog.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService, cfg *config.Config, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
		logger:      logger,
	}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var input domain.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	user, tokens, err := h.authService.Login(c.Request.Context(), input, ipAddress, userAgent)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			response.Unauthorized(c, "invalid email or password")
		case errors.Is(err, domain.ErrAccountLocked):
			response.Unauthorized(c, "account is temporarily locked due to too many failed attempts")
		case errors.Is(err, domain.ErrAccountInactive):
			response.Forbidden(c, "account is inactive")
		default:
			h.logger.Error().Err(err).Msg("login error")
			response.InternalError(c, err)
		}
		return
	}

	// Set refresh token as httpOnly cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken, tokens.ExpiresAt)

	// Return access token in response body
	response.OKWithMessage(c, "login successful", gin.H{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
		"token_type":   tokens.TokenType,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"full_name":  user.FullName,
			"role":       user.Role,
			"avatar_url": user.AvatarURL,
		},
	})
}

// Logout handles POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		// Already logged out or no cookie
		h.clearRefreshTokenCookie(c)
		response.OKWithMessage(c, "logged out successfully", nil)
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		h.logger.Error().Err(err).Msg("logout error")
		// Still clear the cookie even if DB operation fails
	}

	h.clearRefreshTokenCookie(c)
	response.OKWithMessage(c, "logged out successfully", nil)
}

// RefreshToken handles POST /api/v1/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		response.Unauthorized(c, "refresh token not found")
		return
	}

	tokens, err := h.authService.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidToken):
			h.clearRefreshTokenCookie(c)
			response.Unauthorized(c, "invalid or expired refresh token")
		case errors.Is(err, domain.ErrTokenRevoked):
			h.clearRefreshTokenCookie(c)
			response.Unauthorized(c, "refresh token has been revoked")
		case errors.Is(err, domain.ErrAccountInactive):
			h.clearRefreshTokenCookie(c)
			response.Forbidden(c, "account is inactive")
		default:
			h.logger.Error().Err(err).Msg("refresh token error")
			response.InternalError(c, err)
		}
		return
	}

	// Set new refresh token cookie
	h.setRefreshTokenCookie(c, tokens.RefreshToken, tokens.ExpiresAt)

	response.OK(c, gin.H{
		"access_token": tokens.AccessToken,
		"expires_at":   tokens.ExpiresAt,
		"token_type":   tokens.TokenType,
	})
}

// ChangePassword handles POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "authentication required")
		return
	}

	var input domain.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID, input); err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			response.BadRequest(c, "current password is incorrect")
		default:
			h.logger.Error().Err(err).Msg("change password error")
			response.InternalError(c, err)
		}
		return
	}

	// Clear refresh token cookie (force re-login)
	h.clearRefreshTokenCookie(c)
	response.OKWithMessage(c, "password changed successfully", nil)
}

// Me handles GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, _ := c.Get(middleware.ContextKeyUserID)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "authentication required")
		return
	}

	email, _ := c.Get(middleware.ContextKeyEmail)
	role, _ := c.Get(middleware.ContextKeyRole)

	response.OK(c, gin.H{
		"id":    userID,
		"email": email,
		"role":  role,
	})
}

// setRefreshTokenCookie sets the refresh token as a secure httpOnly cookie
func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	sameSite := http.SameSiteStrictMode
	if h.cfg.Cookie.SameSite == "lax" {
		sameSite = http.SameSiteLaxMode
	} else if h.cfg.Cookie.SameSite == "none" {
		sameSite = http.SameSiteNoneMode
	}

	c.SetSameSite(sameSite)
	c.SetCookie(
		refreshTokenCookieName,
		token,
		maxAge,
		"/api/v1/auth",
		h.cfg.Cookie.Domain,
		h.cfg.Cookie.Secure,
		true, // httpOnly
	)
}

// clearRefreshTokenCookie removes the refresh token cookie
func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie(
		refreshTokenCookieName,
		"",
		-1,
		"/api/v1/auth",
		h.cfg.Cookie.Domain,
		h.cfg.Cookie.Secure,
		true,
	)
}
