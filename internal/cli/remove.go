package cli

import "fmt"

type rmCmd struct{}

func (rmCmd) Name() string  { return "rm" }
func (rmCmd) Usage() string { return "rm <id|query> --yes  — delete a task" }

func (rmCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("rm")
	yes := fs.Bool("yes", false, "confirm deletion")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag
	if len(pos) != 1 {
		return fmt.Errorf("rm: exactly one id or query required")
	}
	id, err := resolveArg(cx, pos[0])
	if err != nil {
		return err
	}
	if !*yes {
		return fmt.Errorf("rm: refusing to delete %s without --yes", id)
	}
	if err := cx.Svc.Delete(id); err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Deleted %s", id))
}
