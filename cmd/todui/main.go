// Command todui is a tiny, config-driven task manager: a CLI-first core with a
// Charm Bubble Tea TUI layered on top.
//
// This file is the composition root — the only place that constructs concrete
// adapters and wires them to the interface seams in internal/ports.
package main

import (
	"fmt"
	"os"

	"github.com/develonaut/todui/internal/app"
	"github.com/develonaut/todui/internal/cli"
	"github.com/develonaut/todui/internal/clock"
	tomlcodec "github.com/develonaut/todui/internal/codec/toml"
	"github.com/develonaut/todui/internal/config"
	"github.com/develonaut/todui/internal/ports"
	mdrender "github.com/develonaut/todui/internal/render/markdown"
	filestore "github.com/develonaut/todui/internal/store/file"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	opts, rest := cli.ParseGlobal(args)

	cfg, created, err := config.EnsureConfig(opts.ConfigPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "todui:", err)
		return 1
	}
	if opts.StorePath != "" {
		cfg.Store = opts.StorePath
	}
	if created {
		fmt.Fprintln(os.Stderr, "todui: wrote default config to", opts.ConfigPath)
	}

	renderer := mdrender.New(cfg.Schema())
	svc := buildService(cfg, renderer)

	cx := &cli.Context{
		Svc:      svc,
		Cfg:      cfg,
		Renderer: renderer,
		Out:      os.Stdout,
		Err:      os.Stderr,
	}
	return cli.NewRegistry().Dispatch(cx, rest)
}

// buildService wires the file-backed repository and use-case service.
func buildService(cfg config.Config, renderer ports.Renderer) *app.Service {
	store := filestore.Options{Path: cfg.Store, Codec: tomlcodec.Codec{}}
	if cfg.Mirror != "" {
		store.Renderer = renderer
		store.Mirror = cfg.Mirror
	}
	repo := filestore.New(store)
	set := app.Settings{
		Schema:       cfg.Schema(),
		StartSection: cfg.StartSection,
		DoneMax:      cfg.Done.MaxItems,
		DoneAge:      cfg.Done.MaxAgeDays,
	}
	return app.NewService(repo, clock.System{}, set)
}
