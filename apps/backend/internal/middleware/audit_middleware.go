package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
)

// AuditLogger creates a middleware that automatically logs admin actions to the audit log
func AuditLogger(compRepo repository.ComponentRepository, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only audit write operations
		method := c.Request.Method
		if method == "GET" || method == "OPTIONS" || method == "HEAD" {
			c.Next()
			return
		}

		// Capture request body for audit
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Process request
		c.Next()

		// Only log successful operations (2xx)
		statusCode := c.Writer.Status()
		if statusCode < 200 || statusCode >= 300 {
			return
		}

		// Extract user info from context
		userIDVal, _ := c.Get(ContextKeyUserID)
		emailVal, _ := c.Get(ContextKeyEmail)
		roleVal, _ := c.Get(ContextKeyRole)

		userID, _ := userIDVal.(uuid.UUID)
		email, _ := emailVal.(string)
		role, _ := roleVal.(domain.UserRole)

		// Determine action and resource type from method and path
		action, resourceType, resourceName := extractAuditInfo(method, c.FullPath(), c.Param("id"))

		if action == "" || resourceType == "" {
			return
		}

		// Get site ID if available
		var siteID *uuid.UUID
		siteIDStr := c.Query("site_id")
		if siteIDStr != "" {
			if id, err := uuid.Parse(siteIDStr); err == nil {
				siteID = &id
			}
		}

		// Parse request body for new values
		var newValues domain.JSONMap
		if len(requestBody) > 0 {
			json.Unmarshal(requestBody, &newValues)
			// Remove sensitive fields
			delete(newValues, "password")
			delete(newValues, "password_hash")
			delete(newValues, "current_password")
			delete(newValues, "new_password")
		}

		// Get resource ID from path param
		var resourceID *uuid.UUID
		if idStr := c.Param("id"); idStr != "" {
			if id, err := uuid.Parse(idStr); err == nil {
				resourceID = &id
			}
		}

		// Create audit log entry
		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		roleStr := string(role)

		log := &domain.AuditLog{
			ID:           uuid.New(),
			Action:       action,
			ResourceType: resourceType,
			ResourceID:   resourceID,
			ResourceName: &resourceName,
			NewValues:    newValues,
			IPAddress:    &ipAddress,
			UserAgent:    &userAgent,
			SiteID:       siteID,
		}

		if userID != uuid.Nil {
			log.UserID = &userID
			log.UserEmail = &email
			log.UserRole = &roleStr
		}

		// Save audit log asynchronously to not block response
		// Use background context since request context may be cancelled
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := compRepo.CreateAuditLog(ctx, log); err != nil {
				logger.Error().Err(err).Msg("failed to create audit log")
			}
		}()
	}
}

// extractAuditInfo determines the action, resource type, and resource name from the request
func extractAuditInfo(method, path, id string) (action, resourceType, resourceName string) {
	// Map HTTP methods to actions
	switch method {
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		return "", "", ""
	}

	// Handle special cases
	if strings.Contains(path, "/publish") {
		action = "publish"
	} else if strings.Contains(path, "/unpublish") {
		action = "unpublish"
	} else if strings.Contains(path, "/upload") {
		action = "upload"
	} else if strings.Contains(path, "/change-password") {
		action = "password_change"
	} else if strings.Contains(path, "/login") {
		action = "login"
	} else if strings.Contains(path, "/logout") {
		action = "logout"
	}

	// Extract resource type from path
	pathParts := strings.Split(strings.TrimPrefix(path, "/api/v1/admin/"), "/")
	if len(pathParts) > 0 {
		resourceType = pathParts[0]
		// Normalize resource type
		switch resourceType {
		case "pages":
			resourceType = "page"
		case "sections":
			resourceType = "section"
		case "contents":
			resourceType = "content"
		case "sites":
			resourceType = "site"
		case "users":
			resourceType = "user"
		case "features":
			resourceType = "feature"
		case "testimonials":
			resourceType = "testimonial"
		case "pricing":
			resourceType = "pricing_plan"
		case "faqs":
			resourceType = "faq"
		case "navigation":
			resourceType = "navigation"
		case "media":
			resourceType = "media"
		}
	}

	resourceName = resourceType
	if id != "" {
		resourceName = resourceType + ":" + id
	}

	return action, resourceType, resourceName
}
