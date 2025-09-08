package web

import (
	"embed"

	"github.com/benbjohnson/hashfs"
)

//go:embed static
var StaticDirectory embed.FS

var (
	StaticSys = hashfs.NewFS(StaticDirectory)
)

func StaticPath(path string) string {
	return "/monitor/" + StaticSys.HashName("static/"+path)
}
