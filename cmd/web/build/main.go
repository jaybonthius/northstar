package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	watch := flag.Bool("watch", false, "Enable watcher mode")
	flag.Parse()

	if err := run(*watch); err != nil {
		slog.Error("failure", "error", err)
		os.Exit(1)
	}
}

func findEntryPoints() ([]api.EntryPoint, error) {
	var entryPoints []api.EntryPoint

	features, err := filepath.Glob("app/features/*")
	if err != nil {
		return nil, err
	}

	// filter out non-directories
	var dirs []string
	for _, feature := range features {
		if info, err := os.Stat(feature); err == nil && info.IsDir() {
			dirs = append(dirs, feature)
		}
	}
	features = dirs

	for _, feature := range features {
		featureName := filepath.Base(feature)

		webComponents, err := filepath.Glob(filepath.Join(feature, "web-components/*.ts"))
		if err != nil {
			return nil, err
		}
		for _, wc := range webComponents {
			name := strings.TrimSuffix(filepath.Base(wc), ".ts")
			entryPoints = append(entryPoints, api.EntryPoint{
				InputPath:  wc,
				OutputPath: filepath.Join("app/features", featureName, "web/static/web-components", name),
			})
		}

		styles, err := filepath.Glob(filepath.Join(feature, "styles/*.css"))
		if err != nil {
			return nil, err
		}
		for _, style := range styles {
			name := strings.TrimSuffix(filepath.Base(style), ".css")
			entryPoints = append(entryPoints, api.EntryPoint{
				InputPath:  style,
				OutputPath: filepath.Join("app/features", featureName, "web/static/styles", name),
			})
		}
	}

	return entryPoints, nil
}

func run(watch bool) error {
	entryPoints, err := findEntryPoints()
	if err != nil {
		return err
	}

	opts := api.BuildOptions{
		EntryPointsAdvanced: entryPoints,
		Outdir:              "./",
		Bundle:              true,
		Write:               true,
		LogLevel:            api.LogLevelInfo,
		MinifyWhitespace:    true,
		MinifyIdentifiers:   true,
		MinifySyntax:        true,
		Format:              api.FormatESModule,
		Sourcemap:           api.SourceMapLinked,
		Target:              api.ESNext,
		NodePaths:           []string{"node_modules"},
	}

	if watch {
		slog.Info("Watching...")
		ctx, err := api.Context(opts)
		if err != nil {
			return err
		}

		if err := ctx.Watch(api.WatchOptions{}); err != nil {
			return err
		}

		<-make(chan struct{})
		return nil
	}

	slog.Info("Building...")

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		errs := make([]error, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = errors.New(err.Text)
		}
		return errors.Join(errs...)
	}

	return nil
}
