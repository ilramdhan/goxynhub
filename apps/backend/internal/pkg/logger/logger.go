package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Setup configures and returns a zerolog logger
func Setup(level, format string) zerolog.Logger {
	// Set log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Configure time format
	zerolog.TimeFieldFormat = time.RFC3339

	var logger zerolog.Logger

	if format == "console" || format == "pretty" {
		// Human-readable console output for development
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Caller().Logger()
	} else {
		// JSON output for production
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	return logger
}
