package common

import (
	"northstar/app/features/common/web"
	"northstar/app/static"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	router.Handle("/common/static/*", static.Handler("/common/static", web.StaticDirectory, "common"))
	return nil
}
