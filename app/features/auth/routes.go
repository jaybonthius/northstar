package auth

import (
	"database/sql"

	"northstar/app/features/auth/gen/authdb"
	"northstar/app/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(router chi.Router, db *sql.DB, store sessions.Store) error {
	queries := authdb.New(db)
	authRepository := &authRepository{queries: queries}
	authHandlers := &authHandlers{
		repository: authRepository,
		store:      store,
	}

	router.Route("/login", func(r chi.Router) {
		r.Use(middleware.RedirectIfAuthenticated(store))
		r.Get("/", authHandlers.handleLoginPage)
		r.Post("/", authHandlers.handleLogin)
	})

	router.Route("/signup", func(r chi.Router) {
		r.Use(middleware.RedirectIfAuthenticated(store))
		r.Get("/", authHandlers.handleSignupPage)
		r.Post("/", authHandlers.handleSignup)
	})

	router.Route("/logout", func(r chi.Router) {
		r.Use(middleware.RequireAuth(store, db))
		r.Post("/", authHandlers.handleLogout)
	})

	router.Route("/profile", func(r chi.Router) {
		r.Use(middleware.RequireAuth(store, db))
		r.Get("/", authHandlers.handleProfilePage)
	})

	return nil
}
