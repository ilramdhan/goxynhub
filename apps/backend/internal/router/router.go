package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/config"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/handler"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/middleware"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	AuthHandler *handler.AuthHandler
	PageHandler *handler.PageHandler
	JWTManager  *auth.JWTManager
	Config      *config.Config
	Logger      zerolog.Logger
}

// Setup configures and returns the Gin router
func Setup(deps *Dependencies) *gin.Engine {
	if deps.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.MaxBodySize(deps.Config.Security.MaxUploadSize))

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     deps.Config.CORS.Origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: deps.Config.CORS.AllowCredentials,
		MaxAge:           12 * 3600,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": deps.Config.App.Version,
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")

	// ─── Public Routes ────────────────────────────────────────────────────────
	public := v1.Group("/public")
	if deps.Config.RateLimit.Enabled {
		public.Use(middleware.RateLimiter(deps.Config.RateLimit.Requests))
	}
	{
		// Pages
		public.GET("/pages/:slug", deps.PageHandler.GetPublicPage)
		public.GET("/pages", deps.PageHandler.GetPublicPage) // homepage (no slug)
	}

	// ─── Auth Routes ──────────────────────────────────────────────────────────
	authGroup := v1.Group("/auth")
	if deps.Config.RateLimit.Enabled {
		authGroup.Use(middleware.RateLimiter(deps.Config.RateLimit.AuthRequests))
	}
	{
		authGroup.POST("/login", deps.AuthHandler.Login)
		authGroup.POST("/logout", deps.AuthHandler.Logout)
		authGroup.POST("/refresh", deps.AuthHandler.RefreshToken)

		// Protected auth routes
		authProtected := authGroup.Group("")
		authProtected.Use(middleware.AuthMiddleware(deps.JWTManager))
		{
			authProtected.GET("/me", deps.AuthHandler.Me)
			authProtected.POST("/change-password", deps.AuthHandler.ChangePassword)
		}
	}

	// ─── Admin Routes ─────────────────────────────────────────────────────────
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(deps.JWTManager))
	if deps.Config.RateLimit.Enabled {
		admin.Use(middleware.RateLimiter(200)) // Higher limit for admin
	}
	{
		// Pages (Editor and above)
		pages := admin.Group("/pages")
		pages.Use(middleware.RequireRole(domain.RoleEditor))
		{
			pages.GET("", deps.PageHandler.ListPages)
			pages.POST("", deps.PageHandler.CreatePage)
			pages.GET("/:id", deps.PageHandler.GetPage)
			pages.PUT("/:id", deps.PageHandler.UpdatePage)
			pages.DELETE("/:id", middleware.RequireRole(domain.RoleAdmin), deps.PageHandler.DeletePage)
			pages.PATCH("/:id/publish", deps.PageHandler.PublishPage)
			pages.PATCH("/:id/unpublish", deps.PageHandler.UnpublishPage)
			pages.GET("/:id/sections", deps.PageHandler.ListSections)
			pages.POST("/:id/sections", deps.PageHandler.CreateSection)
		}

		// Sections (Editor and above)
		sections := admin.Group("/sections")
		sections.Use(middleware.RequireRole(domain.RoleEditor))
		{
			sections.PUT("/:id", deps.PageHandler.UpdateSection)
			sections.DELETE("/:id", deps.PageHandler.DeleteSection)
			sections.PATCH("/reorder", deps.PageHandler.ReorderSections)
			sections.GET("/:id/contents", deps.PageHandler.ListContents)
			sections.POST("/:id/contents", deps.PageHandler.UpsertContent)
			sections.POST("/:id/contents/bulk", deps.PageHandler.BulkUpsertContents)
		}

		// Contents (Editor and above)
		contents := admin.Group("/contents")
		contents.Use(middleware.RequireRole(domain.RoleEditor))
		{
			contents.DELETE("/:id", deps.PageHandler.DeleteContent)
		}
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "route not found",
		})
	})

	return r
}
