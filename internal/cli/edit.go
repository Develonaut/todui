package cli

import (
	"flag"
	"fmt"

	"github.com/develonaut/todui/internal/todo"
)

type editCmd struct{}

func (editCmd) Name() string { return "edit" }
func (editCmd) Usage() string {
	return "edit <id|query> [--task T] [--context C] [--ado REF] [--tag T]... [--clear-tags]"
}

func (editCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("edit")
	task := fs.String("task", "", "new task text")
	context := fs.String("context", "", "new context")
	ado := fs.String("ado", "", "new leading reference token")
	clearTags := fs.Bool("clear-tags", false, "remove all tags first")
	var tags multiFlag
	fs.Var(&tags, "tag", "add a tag (repeatable)")
	pos, err := parseArgs(fs, args)
	if err != nil {
		return err
	}
	cx.JSON = *jsonFlag
	if len(pos) < 1 {
		return fmt.Errorf("edit: id or query required")
	}
	id, err := resolveArg(cx, pos[0])
	if err != nil {
		return err
	}

	set := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })

	err = cx.Svc.Edit(id, func(it *todo.Item) {
		if set["task"] {
			it.Task = *task
		}
		if set["context"] {
			it.Context = *context
		}
		if set["ado"] {
			it.ADO = *ado
		}
		if *clearTags {
			it.Tags = nil
		}
		it.Tags = append(it.Tags, tags...)
	})
	if err != nil {
		return err
	}
	return cx.report([]string{id}, fmt.Sprintf("Edited %s", id))
}
