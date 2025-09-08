package sortable

import (
	"net/http"
	"northstar/app/features/sortable/pages"
	"northstar/app/features/sortable/web"
	"northstar/app/static"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(router chi.Router) error {
	router.Handle("/sortable/static/*", static.Handler("/sortable/static", web.StaticDirectory, "sortable"))
	router.Get("/sortable", func(w http.ResponseWriter, r *http.Request) {
		if err := pages.SortablePage().Render(r.Context(), w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	})

	return nil
}
