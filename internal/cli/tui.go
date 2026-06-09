package cli

import (
	"github.com/develonaut/todui/internal/keymap"
	"github.com/develonaut/todui/internal/tui"
)

type tuiCmd struct{}

func (tuiCmd) Name() string  { return "tui" }
func (tuiCmd) Usage() string { return "tui  — launch the interactive UI (default)" }

func (tuiCmd) Run(cx *Context, _ []string) error {
	return tui.Run(cx.Svc, keymap.OverridesFromMap(cx.Cfg.Keys))
}
