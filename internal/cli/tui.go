package cli

import "fmt"

// tuiCmd is a placeholder until the interactive UI lands in a later milestone.
type tuiCmd struct{}

func (tuiCmd) Name() string  { return "tui" }
func (tuiCmd) Usage() string { return "tui  — launch the interactive UI" }

func (tuiCmd) Run(cx *Context, _ []string) error {
	fmt.Fprintln(cx.Err, "todui: the interactive TUI is not built yet — use the CLI (todui help)")
	return nil
}
