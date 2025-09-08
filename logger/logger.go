package logger

import (
	"log/slog"
	"northstar/config"
)

func parseLogLevel() slog.Level {
	logLvlStr := config.Global.LogLevel
	switch logLvlStr {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
