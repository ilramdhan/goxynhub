package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/config"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/domain"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/handler"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/middleware"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
)

// Dependencies holds all handler dependencies
type Dependencies struct {
	AuthHandler      *handler.AuthHandler
	PageHandler      *handler.PageHandler
	SiteHandler      *handler.SiteHandler
	UserHandler      *handler.UserHandler
	ComponentHandler *handler.ComponentHandler
	JWTManager       *auth.JWTManager
	Config           *config.Config
	Logger           zerolog.Logger
	DB               *sqlx.DB // for health check
}

// Setup configures and returns the Gin router
func Setup(deps *Dependencies) *gin.Engine {
	if deps.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.RequestID())
	r.Use(middleware.StructuredLogger(deps.Logger))
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

	// Health check with DB ping
	r.GET("/health", func(c *gin.Context) {
		status := "ok"
		dbStatus := "ok"

		if deps.DB != nil {
			if err := deps.DB.PingContext(c.Request.Context()); err != nil {
				status = "degraded"
				dbStatus = "error"
			}
		}

		statusCode := http.StatusOK
		if status != "ok" {
			statusCode = http.StatusServiceUnavailable
		}

		c.JSON(statusCode, gin.H{
			"status":  status,
			"version": deps.Config.App.Version,
			"env":     deps.Config.App.Env,
			"checks": gin.H{
				"database": dbStatus,
			},
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
		// Site info (for frontend to get site settings, SEO, etc.)
		public.GET("/sites/:id", deps.SiteHandler.GetPublicSiteByID)
		public.GET("/site/:slug", deps.SiteHandler.GetPublicSite)

		// Pages
		public.GET("/pages/:slug", deps.PageHandler.GetPublicPage)
		public.GET("/pages", deps.PageHandler.GetPublicPage) // homepage

		// Navigation
		public.GET("/navigation/:siteId/:identifier", deps.ComponentHandler.GetNavigation)
		public.GET("/navigation/:siteId", deps.ComponentHandler.GetNavigation)

		// Public component endpoints (no auth required for landing page rendering)
		public.GET("/features", deps.ComponentHandler.ListFeatures)
		public.GET("/testimonials", deps.ComponentHandler.ListTestimonials)
		public.GET("/pricing", deps.ComponentHandler.ListPricingPlans)
		public.GET("/faqs", deps.ComponentHandler.ListFAQs)
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
		admin.Use(middleware.RateLimiter(200))
	}
	{
		// ── Sites (Admin+) ──────────────────────────────────────────────────
		sites := admin.Group("/sites")
		sites.Use(middleware.RequireRole(domain.RoleAdmin))
		{
			sites.GET("", deps.SiteHandler.ListSites)
			sites.POST("", deps.SiteHandler.CreateSite)
			sites.GET("/:id", deps.SiteHandler.GetSite)
			sites.PUT("/:id", deps.SiteHandler.UpdateSite)
			sites.DELETE("/:id", middleware.RequireRole(domain.RoleSuperAdmin), deps.SiteHandler.DeleteSite)
			sites.GET("/:id/settings", deps.SiteHandler.GetSettings)
			sites.PUT("/:id/settings", deps.SiteHandler.BulkUpdateSettings)
			sites.PUT("/:id/settings/:key", deps.SiteHandler.UpdateSetting)
		}

		// ── Pages (Editor+) ─────────────────────────────────────────────────
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

		// ── Sections (Editor+) ──────────────────────────────────────────────
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

		// ── Contents (Editor+) ──────────────────────────────────────────────
		contents := admin.Group("/contents")
		contents.Use(middleware.RequireRole(domain.RoleEditor))
		{
			contents.DELETE("/:id", deps.PageHandler.DeleteContent)
		}

		// ── Features (Editor+) ──────────────────────────────────────────────
		features := admin.Group("/features")
		features.Use(middleware.RequireRole(domain.RoleEditor))
		{
			features.GET("", deps.ComponentHandler.ListFeatures)
			features.POST("", deps.ComponentHandler.CreateFeature)
			features.PUT("/:id", deps.ComponentHandler.UpdateFeature)
			features.DELETE("/:id", deps.ComponentHandler.DeleteFeature)
		}

		// ── Testimonials (Editor+) ──────────────────────────────────────────
		testimonials := admin.Group("/testimonials")
		testimonials.Use(middleware.RequireRole(domain.RoleEditor))
		{
			testimonials.GET("", deps.ComponentHandler.ListTestimonials)
			testimonials.POST("", deps.ComponentHandler.CreateTestimonial)
			testimonials.PUT("/:id", deps.ComponentHandler.UpdateTestimonial)
			testimonials.DELETE("/:id", deps.ComponentHandler.DeleteTestimonial)
		}

		// ── Pricing (Editor+) ───────────────────────────────────────────────
		pricing := admin.Group("/pricing")
		pricing.Use(middleware.RequireRole(domain.RoleEditor))
		{
			pricing.GET("", deps.ComponentHandler.ListPricingPlans)
			pricing.POST("", deps.ComponentHandler.CreatePricingPlan)
			pricing.PUT("/:id", deps.ComponentHandler.UpdatePricingPlan)
			pricing.DELETE("/:id", deps.ComponentHandler.DeletePricingPlan)
		}

		// ── FAQs (Editor+) ──────────────────────────────────────────────────
		faqs := admin.Group("/faqs")
		faqs.Use(middleware.RequireRole(domain.RoleEditor))
		{
			faqs.GET("", deps.ComponentHandler.ListFAQs)
			faqs.POST("", deps.ComponentHandler.CreateFAQ)
			faqs.PUT("/:id", deps.ComponentHandler.UpdateFAQ)
			faqs.DELETE("/:id", deps.ComponentHandler.DeleteFAQ)
		}

		// ── Navigation (Editor+) ────────────────────────────────────────────
		navigation := admin.Group("/navigation")
		navigation.Use(middleware.RequireRole(domain.RoleEditor))
		{
			navigation.GET("", deps.ComponentHandler.ListNavigation)
			navigation.POST("/:menuId/items", deps.ComponentHandler.CreateNavigationItem)
			navigation.PUT("/items/:id", deps.ComponentHandler.UpdateNavigationItem)
			navigation.DELETE("/items/:id", deps.ComponentHandler.DeleteNavigationItem)
		}

		// ── Media (Editor+) ─────────────────────────────────────────────────
		media := admin.Group("/media")
		media.Use(middleware.RequireRole(domain.RoleEditor))
		{
			media.GET("", deps.ComponentHandler.ListMedia)
			media.POST("/upload", deps.ComponentHandler.UploadMedia)
			media.PUT("/:id", deps.ComponentHandler.UpdateMedia)
			media.DELETE("/:id", deps.ComponentHandler.DeleteMedia)
		}

		// ── Users (Admin+) ──────────────────────────────────────────────────
		users := admin.Group("/users")
		users.Use(middleware.RequireRole(domain.RoleAdmin))
		{
			users.GET("", deps.UserHandler.ListUsers)
			users.POST("", deps.UserHandler.CreateUser)
			users.GET("/:id", deps.UserHandler.GetUser)
			users.PUT("/:id", deps.UserHandler.UpdateUser)
			users.DELETE("/:id", middleware.RequireRole(domain.RoleSuperAdmin), deps.UserHandler.DeleteUser)
		}

		// ── Audit Logs (Admin+) ─────────────────────────────────────────────
		auditLogs := admin.Group("/audit-logs")
		auditLogs.Use(middleware.RequireRole(domain.RoleAdmin))
		{
			auditLogs.GET("", deps.ComponentHandler.ListAuditLogs)
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
