// Package markdown renders a todo.List as a Markdown document — the read-only
// mirror — reproducing positional IDs, tags, and the section layout.
package markdown

import (
	"strings"

	"github.com/develonaut/todui/internal/ports"
	"github.com/develonaut/todui/internal/todo"
)

// Renderer formats lists using a fixed section schema.
type Renderer struct {
	schema todo.Schema
}

// New returns a Markdown renderer for the given section schema.
func New(schema todo.Schema) *Renderer {
	return &Renderer{schema: schema}
}

var _ ports.Renderer = (*Renderer)(nil)

// Render returns the Markdown representation of the list: header lines verbatim,
// the regenerated "Last updated" line, then one section per schema entry.
func (r *Renderer) Render(l todo.List) ([]byte, error) {
	var b strings.Builder
	for _, h := range l.Header {
		b.WriteString(h)
		b.WriteByte('\n')
	}
	if l.LastUpdated != "" {
		b.WriteString("_Last updated: ")
		b.WriteString(l.LastUpdated)
		b.WriteString("_\n")
	}
	for _, sec := range r.schema.Sections {
		b.WriteString("\n## ")
		b.WriteString(sec.Title)
		b.WriteByte('\n')
		for _, it := range l.Section(r.schema, sec.Key) {
			b.WriteString(r.renderItem(sec, it))
			b.WriteByte('\n')
		}
	}
	return []byte(b.String()), nil
}

// renderItem formats a single item line for its section.
func (r *Renderer) renderItem(sec todo.Section, it todo.Item) string {
	text := todo.JoinTitle(it.Title, it.Description)
	if sec.Done {
		line := "- [x] " + text
		if it.DoneDate != "" {
			line += " (done " + it.DoneDate + ")"
		}
		return line
	}

	var b strings.Builder
	b.WriteString("- [ ] **")
	b.WriteString(r.schema.ID(sec, it.Order))
	b.WriteString("** ")
	if it.ADO != "" {
		b.WriteString(it.ADO)
		b.WriteByte(' ')
	}
	b.WriteString(text)
	for _, tag := range it.Tags {
		b.WriteString(" `[")
		b.WriteString(tag)
		b.WriteString("]`")
	}
	return b.String()
}
