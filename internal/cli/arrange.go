package cli

import (
	"fmt"
	"strings"
)

type mvCmd struct{}

func (mvCmd) Name() string  { return "mv" }
func (mvCmd) Usage() string { return "mv <id|query> <section>  — move to another section" }

func (mvCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("mv")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag
	if len(pos) != 2 {
		return fmt.Errorf("mv: usage: mv <id|query> <section>")
	}
	id, err := resolveArg(cx, pos[0])
	if err != nil {
		return err
	}
	key, ok := resolveSection(cx, pos[1])
	if !ok {
		return fmt.Errorf("mv: unknown section %q", fs.Arg(1))
	}
	if err := cx.Svc.Move(id, key); err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Moved %s to %s", id, key))
}

type reorderCmd struct{}

func (reorderCmd) Name() string  { return "reorder" }
func (reorderCmd) Usage() string { return "reorder <id|query> up|down  — reorder within a section" }

func (reorderCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("reorder")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag
	if len(pos) != 2 {
		return fmt.Errorf("reorder: usage: reorder <id|query> up|down")
	}
	id, err := resolveArg(cx, pos[0])
	if err != nil {
		return err
	}
	var step int
	dir := strings.ToLower(pos[1])
	switch dir {
	case "up":
		step = -1
	case "down":
		step = 1
	default:
		return fmt.Errorf("reorder: direction must be up or down")
	}
	if err := cx.Svc.Reorder(id, step); err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Reordered %s %s", id, dir))
}
