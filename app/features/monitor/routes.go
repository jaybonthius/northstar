package monitor

import (
	"northstar/app/features/monitor/web"
	"northstar/app/static"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	handlers := NewHandlers()

	router.Handle("/monitor/static/*", static.Handler("/monitor/static", web.StaticDirectory, "monitor"))
	router.Get("/monitor", handlers.MonitorPage)
	router.Get("/monitor/events", handlers.MonitorEvents)

	return nil
}
