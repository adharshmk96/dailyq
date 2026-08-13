package logger

import (
	"log/slog"
	"os"
	"strings"

	"dailyq-api/internal/config"
)

// New builds a slog logger from the log configuration.
func New(cfg config.LogConfig) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(cfg.Level)}

	var handler slog.Handler
	if strings.EqualFold(cfg.Format, "json") {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// SetDefault installs the logger as the process-wide default.
func SetDefault(l *slog.Logger) {
	slog.SetDefault(l)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
