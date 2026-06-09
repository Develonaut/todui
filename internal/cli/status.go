package cli

import "fmt"

type doneCmd struct{}

func (doneCmd) Name() string  { return "done" }
func (doneCmd) Usage() string { return "done <id|query> [--json]  — complete a task" }

func (doneCmd) Run(cx *Context, args []string) error {
	id, err := singleArg(cx, "done", args)
	if err != nil {
		return err
	}
	if err := cx.Svc.Complete(id); err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Completed %s", id))
}

type startCmd struct{}

func (startCmd) Name() string  { return "start" }
func (startCmd) Usage() string { return "start <id|query> [--json]  — claim/begin a task" }

func (startCmd) Run(cx *Context, args []string) error {
	id, err := singleArg(cx, "start", args)
	if err != nil {
		return err
	}
	if err := cx.Svc.Start(id); err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Started %s", id))
}

// singleArg parses the shared --json flag and resolves the one positional
// id-or-query argument these commands take.
func singleArg(cx *Context, name string, args []string) (string, error) {
	fs, jsonFlag := newFlagSet(name)
	pos, err := parseArgs(fs, args)
	if err != nil {
		return "", err
	}
	cx.JSON = *jsonFlag
	if len(pos) != 1 {
		return "", fmt.Errorf("%s: exactly one id or query required", name)
	}
	return resolveArg(cx, pos[0])
}
