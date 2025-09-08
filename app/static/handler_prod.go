//go:build !dev
// +build !dev

package static

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/benbjohnson/hashfs"
)

var staticSystems = make(map[string]*hashfs.FS)

func Handler(prefix string, embeddedFS fs.FS, staticPath string) http.Handler {
	slog.Debug("static assets are embedded", "staticPath", staticPath)
	staticSys := hashfs.NewFS(embeddedFS)
	staticSystems[staticPath] = staticSys

	featurePrefix := "/" + staticPath
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix(featurePrefix, hashfs.FileServer(staticSys)).ServeHTTP(w, r)
	})
}

func StaticPath(featurePrefix, path string) string {
	staticSys := staticSystems[featurePrefix]
	if staticSys == nil {
		return "/" + featurePrefix + "/static/" + path
	}
	return "/" + featurePrefix + "/" + staticSys.HashName("static/"+path)
}
