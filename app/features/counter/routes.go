package counter

import (
	"northstar/app/features/counter/web"
	"northstar/app/static"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(router chi.Router, sessionStore sessions.Store) error {
	handlers := NewHandlers(sessionStore)

	router.Handle("/counter/static/*", static.Handler("/counter/static", web.StaticDirectory, "counter"))
	router.Get("/counter", handlers.CounterPage)
	router.Get("/counter/data", handlers.CounterData)

	router.Route("/counter/increment", func(incrementRouter chi.Router) {
		incrementRouter.Post("/global", handlers.IncrementGlobal)
		incrementRouter.Post("/user", handlers.IncrementUser)
	})

	return nil
}
