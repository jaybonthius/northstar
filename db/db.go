package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func InitDatabase() (*sql.DB, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	database, err := sql.Open("sqlite", "data/northstar.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = database.Ping(); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			slog.Error("Failed to close database", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err = goose.SetDialect("sqlite"); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			slog.Error("Failed to close database", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	slog.Info("running database migrations")
	goose.SetBaseFS(MigrationFiles)
	if err = goose.Up(database, "migrations"); err != nil {
		if closeErr := database.Close(); closeErr != nil {
			slog.Error("Failed to close database", "error", closeErr)
		}
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("database initialized successfully")
	return database, nil
}
