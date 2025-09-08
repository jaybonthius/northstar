//go:build dev

package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func CreateLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:   parseLogLevel(),
		NoColor: !isatty.IsTerminal(os.Stdout.Fd()),
	}))
}
