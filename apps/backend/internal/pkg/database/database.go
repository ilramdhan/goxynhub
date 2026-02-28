package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/ilramdhan/goxynhub/apps/backend/internal/config"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("database.Connect open: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database.Connect ping: %w", err)
	}

	return db, nil
}
