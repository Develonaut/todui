package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/develonaut/todui/internal/todo"
)

// row is an item plus its computed display ID, for JSON output.
type row struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	ADO         string   `json:"ado,omitempty"`
	Section     string   `json:"section"`
	Done        string   `json:"done_date,omitempty"`
}

// rowsFor builds JSON rows for a normalized list.
func rowsFor(l todo.List, s todo.Schema) []row {
	rows := make([]row, 0, len(l.Items))
	for i := range l.Items {
		it := l.Items[i]
		rows = append(rows, row{
			ID: l.ComputeID(s, i), Title: it.Title, Description: it.Description, Tags: it.Tags,
			ADO: it.ADO, Section: it.Section, Done: it.DoneDate,
		})
	}
	return rows
}

// emitJSON writes a single machine-readable result object.
func emitJSON(w io.Writer, ok bool, ids []string, rows []row, e error) error {
	out := map[string]any{"ok": ok}
	if len(ids) > 0 {
		out["ids"] = ids
	}
	if rows != nil {
		out["items"] = rows
	}
	if e != nil {
		out["error"] = e.Error()
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// filterRows keeps only rows in the given section.
func filterRows(rows []row, section string) []row {
	out := make([]row, 0, len(rows))
	for _, r := range rows {
		if r.Section == section {
			out = append(out, r)
		}
	}
	return out
}

// printList renders a human-readable, section-grouped view. When only is
// non-empty, just that section is shown.
func (cx *Context) printList(l todo.List, only string) {
	s := cx.Svc.Schema()
	for _, sec := range s.Sections {
		if only != "" && sec.Key != only {
			continue
		}
		items := l.Section(s, sec.Key)
		fmt.Fprintf(cx.Out, "%s (%d)\n", sec.Title, len(items))
		for _, it := range items {
			id := s.ID(sec, it.Order)
			mark := id
			if mark == "" {
				mark = "✓"
			}
			fmt.Fprintf(cx.Out, "  %-3s %s", mark, it.Title)
			for _, tag := range it.Tags {
				fmt.Fprintf(cx.Out, " [%s]", tag)
			}
			if it.DoneDate != "" {
				fmt.Fprintf(cx.Out, " (done %s)", it.DoneDate)
			}
			fmt.Fprintln(cx.Out)
		}
	}
}
