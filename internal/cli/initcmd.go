package cli

import "fmt"

type initCmd struct{}

func (initCmd) Name() string  { return "init" }
func (initCmd) Usage() string { return "init  — create the store (config is created automatically)" }

func (initCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("init")
	if _, err := parseArgs(fs, args); err != nil {
		return err
	}
	cx.JSON = *jsonFlag
	if err := cx.Svc.Touch(); err != nil {
		return err
	}
	return cx.report(nil, fmt.Sprintf("Initialized store at %s", cx.Cfg.Store))
}
