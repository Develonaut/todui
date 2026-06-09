package cli

import (
	"fmt"
	"os"
)

type renderCmd struct{}

func (renderCmd) Name() string  { return "render" }
func (renderCmd) Usage() string { return "render [--out FILE]  — print or write the Markdown view" }

func (renderCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("render")
	out := fs.String("out", "", "output file (default: stdout)")
	if _, err := parseArgs(fs, args); err != nil {
		return err
	}
	cx.JSON = *jsonFlag

	l, err := cx.Svc.List()
	if err != nil {
		return err
	}
	md, err := cx.Renderer.Render(l)
	if err != nil {
		return err
	}

	if *out == "" {
		_, err := cx.Out.Write(md)
		return err
	}
	_ = os.Remove(*out) // allow overwriting a read-only mirror
	if err := os.WriteFile(*out, md, 0o644); err != nil {
		return err
	}
	return cx.report(nil, fmt.Sprintf("Wrote %s", *out))
}
