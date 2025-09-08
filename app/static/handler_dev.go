//go:build dev
// +build dev

package static

import (
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func Handler(prefix string, embeddedFS fs.FS, staticPath string) http.Handler {
	basePath := filepath.Join("app", "features", staticPath, "web", "static")
	slog.Info("static assets are being served from filesystem", "path", basePath)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		staticFS := os.DirFS(basePath)
		http.StripPrefix(prefix, http.FileServerFS(staticFS)).ServeHTTP(w, r)
	})
}

func StaticPath(featurePrefix, path string) string {
	return "/" + featurePrefix + "/static/" + path
}
