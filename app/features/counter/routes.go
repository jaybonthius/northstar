package counter

import (
	"net/http"
	"northstar/app/features/counter/web"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(router chi.Router, sessionStore sessions.Store) error {
	handlers := NewHandlers(sessionStore)

	router.Handle("/counter/static/*", http.StripPrefix("/counter", hashfs.FileServer(web.StaticSys)))
	router.Get("/counter", handlers.CounterPage)
	router.Get("/counter/data", handlers.CounterData)

	router.Route("/counter/increment", func(incrementRouter chi.Router) {
		incrementRouter.Post("/global", handlers.IncrementGlobal)
		incrementRouter.Post("/user", handlers.IncrementUser)
	})

	return nil
}
