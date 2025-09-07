package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"

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

func run(watch bool) error {
	opts := api.BuildOptions{
		EntryPointsAdvanced: []api.EntryPoint{
			{
				InputPath:  "app/features/reverse/web-components/reverse-component.ts",
				OutputPath: "app/features/reverse/web/static/web-components/reverse-component",
			},
			{
				InputPath:  "app/features/sortable/web-components/sortable-example.ts",
				OutputPath: "app/features/sortable/web/static/web-components/sortable-example",
			},
			{
				InputPath:  "app/features/common/styles/styles.css",
				OutputPath: "app/features/common/web/static/index",
			},
			{
				InputPath:  "app/features/sortable/styles/styles.css",
				OutputPath: "app/features/sortable/web/static/index",
			},
			{
				InputPath:  "app/features/counter/styles/styles.css",
				OutputPath: "app/features/counter/web/static/index",
			},
			{
				InputPath:  "app/features/index/styles/styles.css",
				OutputPath: "app/features/index/web/static/index",
			},
			{
				InputPath:  "app/features/reverse/styles/styles.css",
				OutputPath: "app/features/reverse/web/static/index",
			},
			{
				InputPath:  "app/features/monitor/styles/styles.css",
				OutputPath: "app/features/monitor/web/static/index",
			},
		},
		Outdir:            "./",
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Format:            api.FormatESModule,
		Sourcemap:         api.SourceMapLinked,
		Target:            api.ESNext,
		NodePaths:         []string{"node_modules"},
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
