package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup initializes the global logger based on configuration
func Setup(level, format string) zerolog.Logger {
	// Set log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Set time format
	zerolog.TimeFieldFormat = time.RFC3339

	var logger zerolog.Logger
	if format == "console" {
		// Human-readable output for development
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Set as global logger
	log.Logger = logger

	return logger
}

// FromContext extracts the logger from context or returns the global logger
func FromContext(ctx interface{ Value(key interface{}) interface{} }) zerolog.Logger {
	if l, ok := ctx.Value("logger").(zerolog.Logger); ok {
		return l
	}
	return log.Logger
}
