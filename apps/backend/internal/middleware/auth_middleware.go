package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/response"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyEmail  = "user_email"
	ContextKeyRole   = "user_role"
	ContextKeyClaims = "jwt_claims"
)

// AuthMiddleware validates JWT access tokens
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// Store claims in context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(roles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextKeyRole)
		if !exists {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		userRole, ok := roleVal.(domain.UserRole)
		if !ok {
			response.Unauthorized(c, "invalid token claims")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			// Create a temporary user to use the HasRole method
			tempUser := &domain.User{Role: userRole}
			if tempUser.HasRole(requiredRole) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractBearerToken extracts the Bearer token from the Authorization header
func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}

	return parts[1]
}
