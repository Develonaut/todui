package cli

import (
	"fmt"
	"strings"

	"github.com/develonaut/todui/internal/todo"
)

type addCmd struct{}

func (addCmd) Name() string { return "add" }
func (addCmd) Usage() string {
	return "add <title...> [--desc D] [--section S] [--tag T]... [--ado REF]"
}

func (addCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("add")
	section := fs.String("section", "", "section key (default: first non-done)")
	desc := fs.String("desc", "", "fuller description")
	ado := fs.String("ado", "", "leading reference token (e.g. #123)")
	var tags multiFlag
	fs.Var(&tags, "tag", "tag slug (repeatable)")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag

	text := strings.TrimSpace(strings.Join(pos, " "))
	if text == "" {
		return fmt.Errorf("add: title required")
	}
	title, description := todo.SplitTitle(text)
	if *desc != "" {
		description = *desc
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
		Title: title, Description: description, ADO: *ado, Tags: tags, Section: sec,
	})
	if err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Added %s — %s", id, title))
}
