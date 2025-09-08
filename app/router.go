package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"northstar/config"
	"northstar/nats"
	"sync"
	"time"

	"northstar/app/features/common"
	"northstar/app/features/counter"
	"northstar/app/features/index"
	"northstar/app/features/monitor"
	"northstar/app/features/reverse"
	"northstar/app/features/sortable"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/starfederation/datastar-go/datastar"
)

func SetupRoutes(ctx context.Context, router chi.Router) (err error) {
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

	ns, err := nats.SetupNATS(ctx)
	if err != nil {
		return fmt.Errorf("error setting up NATS: %w", err)
	}

	sessionStore := sessions.NewCookieStore([]byte("session-secret"))
	sessionStore.MaxAge(int(24 * time.Hour / time.Second))

	if err := errors.Join(
		common.SetupRoutes(router),
		index.SetupRoutes(router, sessionStore, ns),
		counter.SetupRoutes(router, sessionStore),
		monitor.SetupRoutes(router),
		sortable.SetupRoutes(router),
		reverse.SetupRoutes(router),
	); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	return nil
}
