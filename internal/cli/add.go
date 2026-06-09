package cli

import (
	"fmt"
	"strings"

	"github.com/develonaut/todui/internal/todo"
)

type addCmd struct{}

func (addCmd) Name() string { return "add" }
func (addCmd) Usage() string {
	return "add <task...> [--section S] [--tag T]... [--ado REF] [--context C] [--claimed]"
}

func (addCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("add")
	section := fs.String("section", "", "section key (default: first non-done)")
	context := fs.String("context", "", "trailing context after the task")
	ado := fs.String("ado", "", "leading reference token (e.g. #123)")
	claimed := fs.Bool("claimed", false, "mark as claimed")
	var tags multiFlag
	fs.Var(&tags, "tag", "tag slug (repeatable)")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag

	task := strings.TrimSpace(strings.Join(pos, " "))
	if task == "" {
		return fmt.Errorf("add: task text required")
	}

	sec := cx.Cfg.DefaultSection()
	if *section != "" {
		key, ok := resolveSection(cx, *section)
		if !ok {
			return fmt.Errorf("add: unknown section %q", *section)
		}
		sec = key
	}

	id, err := cx.Svc.Add(todo.Item{
		Task: task, Context: *context, ADO: *ado, Tags: tags, Claimed: *claimed, Section: sec,
	})
	if err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Added %s — %s", id, task))
}
