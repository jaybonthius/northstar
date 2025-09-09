package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"northstar/config"
	"sync"

	"northstar/app/middleware"

	"northstar/app/features/auth"
	"northstar/app/features/common"
	"northstar/app/features/counter"
	"northstar/app/features/index"
	"northstar/app/features/monitor"
	"northstar/app/features/reverse"
	"northstar/app/features/sortable"

	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/starfederation/datastar-go/datastar"
)

func SetupRoutes(ctx context.Context, router chi.Router, db *sql.DB, sessionStore *sessions.CookieStore, ns *embeddednats.Server) (err error) {
	// apply optional auth middleware to all routes
	router.Use(middleware.WithAuth(sessionStore, db))

	// setup auth routes
	if err := auth.SetupRoutes(router, db, sessionStore); err != nil {
		return fmt.Errorf("error setting up auth routes: %w", err)
	}

	// setup unprotected routes
	if err := errors.Join(
		common.SetupRoutes(router),
		index.SetupRoutes(router, sessionStore, ns),
		counter.SetupRoutes(router, sessionStore),
		monitor.SetupRoutes(router),
		sortable.SetupRoutes(router),
		reverse.SetupRoutes(router),
	); err != nil {
		return fmt.Errorf("error setting up unprotected routes: %w", err)
	}

	// TODO make some of the routes protected
	// // setup protected routes with auth middleware
	// var protectedRouteErr error
	// router.Group(func(r chi.Router) {
	// 	r.Use(auth.RequireAuth(sessionStore, db))
	// 	protectedRouteErr = errors.Join()
	// })
	// if protectedRouteErr != nil {
	// 	return fmt.Errorf("error setting up protected routes: %w", protectedRouteErr)
	// }

	// setup reload routes
	reloadChan := make(chan struct{}, 1)
	var hotReloadOnce sync.Once
	router.Get("/reload", func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)
		reload := func() { sse.ExecuteScript("window.location.reload()") }
		hotReloadOnce.Do(reload)
		select {
		case <-reloadChan:
			reload()
		case <-r.Context().Done():
		}
	})

	if config.Global.Environment == config.Dev {
		router.Get("/force-reload", func(w http.ResponseWriter, r *http.Request) {
			select {
			case reloadChan <- struct{}{}:
			default:
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	}

	return nil
}
