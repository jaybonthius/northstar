package reverse

import (
	"net/http"
	"northstar/app/features/reverse/pages"
	"northstar/app/features/reverse/web"

	"github.com/benbjohnson/hashfs"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	router.Handle("/reverse/static/*", http.StripPrefix("/reverse", hashfs.FileServer(web.StaticSys)))
	router.Get("/reverse", func(w http.ResponseWriter, r *http.Request) {
		if err := pages.ReversePage().Render(r.Context(), w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	return nil
}
