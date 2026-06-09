package cli

import (
	"fmt"
	"strings"
)

type resolveCmd struct{}

func (resolveCmd) Name() string { return "resolve" }
func (resolveCmd) Usage() string {
	return "resolve <query> [--json]  — list candidate IDs for a query"
}

func (resolveCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("resolve")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag
	if len(pos) < 1 {
		return fmt.Errorf("resolve: query required")
	}
	ids, err := candidates(cx, strings.Join(pos, " "))
	if err != nil {
		return err
	}
	if cx.JSON {
		return emitJSON(cx.Out, true, ids, nil, nil)
	}
	if len(ids) == 0 {
		return fmt.Errorf("no match")
	}
	for _, id := range ids {
		fmt.Fprintln(cx.Out, id)
	}
	return nil
}
