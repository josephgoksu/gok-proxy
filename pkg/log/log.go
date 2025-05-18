package log

import (
	"log/slog"
	"os"
)

func NewLogger(level string) (*slog.Logger, error) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default: // "info" and any other unspecified value
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Default to JSON handler, could be made configurable
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return logger, nil
}
