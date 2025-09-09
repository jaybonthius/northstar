package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	app "northstar/app"
	"northstar/config"
	"northstar/db"
	"northstar/logger"
	"northstar/nats"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := logger.CreateLogger()
	slog.SetDefault(logger)

	if err := run(ctx); err != nil && err != http.ErrServerClosed {
		slog.Error("error running server", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	slog.Info("Configuration loaded", "host", config.Global.Host, "port", config.Global.Port, "log_level", config.Global.LogLevel, "environment", config.Global.Environment)

	// Initialize Database
	database, err := db.InitDatabase()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			slog.Error("error closing database", slog.Any("error", err))
		}
	}()

	// Initialize NATS
	ns, err := nats.SetupNATS(ctx)
	if err != nil {
		return fmt.Errorf("error setting up NATS: %w", err)
	}

	addr := fmt.Sprintf("%s:%s", config.Global.Host, config.Global.Port)
	slog.Info("server started", "host", config.Global.Host, "port", config.Global.Port)
	defer slog.Info("server shutdown complete")

	eg, egctx := errgroup.WithContext(ctx)

	sessionKey := config.Global.SessionSecret
	if sessionKey == "" {
		sessionKey = "dev-session-key-change-in-production-very-long-key"
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false // Set to true in production with HTTPS
	store.Options.SameSite = http.SameSiteLaxMode

	router := chi.NewMux()
	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	if err := app.SetupRoutes(egctx, router, database, store, ns); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return egctx
		},
		ErrorLog: slog.NewLogLogger(
			slog.Default().Handler(),
			slog.LevelError,
		),
	}

	eg.Go(func() error {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		<-egctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Debug("shutting down server...")

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("error during shutdown", "error", err)
			return err
		}

		return nil
	})

	return eg.Wait()
}
