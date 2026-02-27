package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	CORS      CORSConfig
	Supabase  SupabaseConfig
	RateLimit RateLimitConfig
	Log       LogConfig
	Security  SecurityConfig
	Cookie    CookieConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Env     string
	Port    string
	Name    string
	Version string
	Debug   bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	Issuer        string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	Origins          []string
	AllowCredentials bool
}

// SupabaseConfig holds Supabase configuration
type SupabaseConfig struct {
	URL           string
	AnonKey       string
	ServiceKey    string
	StorageBucket string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled      bool
	Requests     int
	Window       time.Duration
	AuthRequests int
	AuthWindow   time.Duration
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	BcryptCost       int
	MaxUploadSize    int64
	AllowedMimeTypes []string
}

// CookieConfig holds cookie configuration
type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite string
}

// Load reads configuration from environment variables and .env file
func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Read .env file (ignore error if not found - use env vars directly)
	_ = viper.ReadInConfig()

	// Set defaults
	setDefaults()

	cfg := &Config{
		App: AppConfig{
			Env:     viper.GetString("APP_ENV"),
			Port:    viper.GetString("APP_PORT"),
			Name:    viper.GetString("APP_NAME"),
			Version: viper.GetString("APP_VERSION"),
			Debug:   viper.GetBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			URL:             viper.GetString("DATABASE_URL"),
			MaxOpenConns:    viper.GetInt("DATABASE_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DATABASE_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DATABASE_CONN_MAX_LIFETIME"),
		},
		JWT: JWTConfig{
			AccessSecret:  viper.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: viper.GetString("JWT_REFRESH_SECRET"),
			AccessExpiry:  viper.GetDuration("JWT_ACCESS_EXPIRY"),
			RefreshExpiry: viper.GetDuration("JWT_REFRESH_EXPIRY"),
			Issuer:        viper.GetString("JWT_ISSUER"),
		},
		CORS: CORSConfig{
			Origins:          strings.Split(viper.GetString("CORS_ORIGINS"), ","),
			AllowCredentials: viper.GetBool("CORS_ALLOW_CREDENTIALS"),
		},
		Supabase: SupabaseConfig{
			URL:           viper.GetString("SUPABASE_URL"),
			AnonKey:       viper.GetString("SUPABASE_ANON_KEY"),
			ServiceKey:    viper.GetString("SUPABASE_SERVICE_KEY"),
			StorageBucket: viper.GetString("SUPABASE_STORAGE_BUCKET"),
		},
		RateLimit: RateLimitConfig{
			Enabled:      viper.GetBool("RATE_LIMIT_ENABLED"),
			Requests:     viper.GetInt("RATE_LIMIT_REQUESTS"),
			Window:       viper.GetDuration("RATE_LIMIT_WINDOW"),
			AuthRequests: viper.GetInt("RATE_LIMIT_AUTH_REQUESTS"),
			AuthWindow:   viper.GetDuration("RATE_LIMIT_AUTH_WINDOW"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		Security: SecurityConfig{
			BcryptCost:       viper.GetInt("BCRYPT_COST"),
			MaxUploadSize:    viper.GetInt64("MAX_UPLOAD_SIZE"),
			AllowedMimeTypes: strings.Split(viper.GetString("ALLOWED_MIME_TYPES"), ","),
		},
		Cookie: CookieConfig{
			Domain:   viper.GetString("COOKIE_DOMAIN"),
			Secure:   viper.GetBool("COOKIE_SECURE"),
			SameSite: viper.GetString("COOKIE_SAME_SITE"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks that required configuration values are set
func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	if c.JWT.AccessSecret == c.JWT.RefreshSecret {
		return fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be different")
	}
	return nil
}

// IsProduction returns true if the app is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// IsDevelopment returns true if the app is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

func setDefaults() {
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_NAME", "landing-cms-api")
	viper.SetDefault("APP_VERSION", "1.0.0")
	viper.SetDefault("APP_DEBUG", false)

	viper.SetDefault("DATABASE_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DATABASE_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DATABASE_CONN_MAX_LIFETIME", "5m")

	viper.SetDefault("JWT_ACCESS_EXPIRY", "15m")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	viper.SetDefault("JWT_ISSUER", "landing-cms-api")

	viper.SetDefault("CORS_ALLOW_CREDENTIALS", true)

	viper.SetDefault("SUPABASE_STORAGE_BUCKET", "media")

	viper.SetDefault("RATE_LIMIT_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_REQUESTS", 100)
	viper.SetDefault("RATE_LIMIT_WINDOW", "1m")
	viper.SetDefault("RATE_LIMIT_AUTH_REQUESTS", 5)
	viper.SetDefault("RATE_LIMIT_AUTH_WINDOW", "1m")

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("BCRYPT_COST", 12)
	viper.SetDefault("MAX_UPLOAD_SIZE", 10485760) // 10MB
	viper.SetDefault("ALLOWED_MIME_TYPES", "image/jpeg,image/png,image/gif,image/webp,image/svg+xml,video/mp4,application/pdf")

	viper.SetDefault("COOKIE_DOMAIN", "localhost")
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("COOKIE_SAME_SITE", "strict")
}
