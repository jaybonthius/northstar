package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	app "northstar/app"
	"northstar/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: func() slog.Level {
			switch config.Global.LogLevel {
			case "DEBUG":
				return slog.LevelDebug
			case "INFO":
				return slog.LevelInfo
			case "WARN":
				return slog.LevelWarn
			case "ERROR":
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		}(),
	}))
	slog.SetDefault(logger)

	if err := run(ctx); err != nil && err != http.ErrServerClosed {
		slog.Error("error running server", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%s", config.Global.Host, config.Global.Port)
	slog.Info("server started", "host", config.Global.Host, "port", config.Global.Port)
	defer slog.Info("server shutdown complete")

	eg, egctx := errgroup.WithContext(ctx)

	router := chi.NewMux()
	router.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	if err := app.SetupRoutes(egctx, router); err != nil {
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
