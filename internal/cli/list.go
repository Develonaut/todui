package cli

import "fmt"

type listCmd struct{}

func (listCmd) Name() string  { return "list" }
func (listCmd) Usage() string { return "list [--section S] [--json]  — show tasks" }

func (listCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("list")
	section := fs.String("section", "", "limit to one section")
	if _, err := parseArgs(fs, args); err != nil {
		return err
	}
	cx.JSON = *jsonFlag

	l, err := cx.Svc.List()
	if err != nil {
		return err
	}

	key := ""
	if *section != "" {
		k, ok := resolveSection(cx, *section)
		if !ok {
			return fmt.Errorf("list: unknown section %q", *section)
		}
		key = k
	}

	if cx.JSON {
		rows := rowsFor(l, cx.Svc.Schema())
		if key != "" {
			rows = filterRows(rows, key)
		}
		return emitJSON(cx.Out, true, nil, rows, nil)
	}
	cx.printList(l, key)
	return nil
}
