package common

import (
	"net/http"
	"northstar/app/features/common/web"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	router.Handle("/common/static/*", http.StripPrefix("/common", hashfs.FileServer(web.StaticSys)))
	return nil
}
