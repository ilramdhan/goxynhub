package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/config"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/handler"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/auth"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/database"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/pkg/logger"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/repository"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/router"
	"github.com/ilramdhan/goxynhub/apps/backend/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Setup logger
	appLogger := logger.Setup(cfg.Log.Level, cfg.Log.Format)
	appLogger.Info().
		Str("env", cfg.App.Env).
		Str("version", cfg.App.Version).
		Msg("starting application")

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()
	appLogger.Info().Msg("database connected successfully")

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
		cfg.JWT.Issuer,
	)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	pageRepo := repository.NewPageRepository(db)

	// Initialize services
	authSvc := service.NewAuthService(userRepo, jwtManager, appLogger)
	pageSvc := service.NewPageService(pageRepo, appLogger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc, cfg, appLogger)
	pageHandler := handler.NewPageHandler(pageSvc, appLogger)

	// Setup router
	deps := &router.Dependencies{
		AuthHandler: authHandler,
		PageHandler: pageHandler,
		JWTManager:  jwtManager,
		Config:      cfg,
		Logger:      appLogger,
	}
	r := router.Setup(deps)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info().
			Str("port", cfg.App.Port).
			Msg("server starting")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info().Msg("shutting down server...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error().Err(err).Msg("server forced to shutdown")
	}

	appLogger.Info().Msg("server exited")
	log.Info().Msg("goodbye!")
}
