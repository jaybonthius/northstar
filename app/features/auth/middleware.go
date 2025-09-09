package auth

import (
	"context"
	"database/sql"
	"net/http"

	"northstar/app/features/auth/gen/authdb"

	"github.com/gorilla/sessions"
)

type contextKey string

const UserContextKey = contextKey("user")

func RequireAuth(store sessions.Store, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if GetUserIDFromContext(r.Context()) == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RedirectIfAuthenticated(store sessions.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "auth-session")

			userID, ok := session.Values["user_id"]
			if ok && userID != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func WithAuth(store sessions.Store, db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			queries := authdb.New(db)
			session, _ := store.Get(r, "auth-session")

			userID, ok := session.Values["user_id"]
			if ok && userID != nil {
				userIDStr := userID.(string)
				user, err := queries.GetUser(r.Context(), userIDStr)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserContextKey, user)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	user, ok := ctx.Value(UserContextKey).(authdb.User)
	if !ok {
		return ""
	}
	return user.ID
}

func IsAuthenticated(r *http.Request, store sessions.Store, db *sql.DB) bool {
	user, _ := GetAuthenticatedUser(r, store, db)
	return user != nil
}

func GetAuthenticatedUser(r *http.Request, store sessions.Store, db *sql.DB) (*authdb.User, error) {
	queries := authdb.New(db)
	session, _ := store.Get(r, "auth-session")

	userID, ok := session.Values["user_id"]
	if !ok || userID == nil {
		return nil, nil
	}

	userIDStr := userID.(string)
	user, err := queries.GetUser(r.Context(), userIDStr)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserAuthStatus(r *http.Request, store sessions.Store, db *sql.DB) bool {
	return IsAuthenticated(r, store, db)
}
