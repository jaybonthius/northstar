//go:build !dev

package logger

import (
	"log/slog"
	"os"
)

func CreateLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(),
	}))
}
