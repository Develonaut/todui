package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"

	"github.com/develonaut/todui/internal/todo"
)

// enterAdd opens an empty form for a new task in the current section.
func (m *Model) enterAdd() tea.Cmd {
	m.editID, m.editSec = "", ""
	m.fTitle, m.fDesc, m.fTags, m.fADO = "", "", "", ""
	m.fSection = m.defaultFormSection()
	m.form = m.buildForm()
	m.mode = modeForm
	return m.form.Init()
}

// enterEdit opens a form pre-filled from the selected item.
func (m *Model) enterEdit() tea.Cmd {
	r, ok := m.selectedRow()
	if !ok || r.id == "" {
		return nil
	}
	m.editID, m.editSec = r.id, r.section.Key
	m.fTitle = r.item.Title
	m.fDesc = r.item.Description
	m.fTags = strings.Join(r.item.Tags, ", ")
	m.fADO = r.item.ADO
	m.fSection = r.section.Key
	m.form = m.buildForm()
	m.mode = modeForm
	return m.form.Init()
}

// exitForm returns to the list and refreshes.
func (m *Model) exitForm() {
	m.mode = modeList
	m.form = nil
	m.rebuild()
}

// buildForm constructs the add/edit form bound to the model's field values.
func (m *Model) buildForm() *huh.Form {
	var opts []huh.Option[string]
	for _, sec := range m.svc.Schema().Sections {
		if !sec.Done {
			opts = append(opts, huh.NewOption(sec.Title, sec.Key))
		}
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Title").Value(&m.fTitle),
			huh.NewText().Title("Description").Value(&m.fDesc),
			huh.NewInput().Title("Tags (comma-separated)").Value(&m.fTags),
			huh.NewInput().Title("Reference").Placeholder("#123").Value(&m.fADO),
			huh.NewSelect[string]().Title("Section").Options(opts...).Value(&m.fSection),
		),
	).WithWidth(formWidth(m.width)).WithHeight(formHeight(m.height))
}

// applyForm writes the form's values back through the service.
func (m *Model) applyForm() {
	tags := splitTags(m.fTags)
	if m.editID == "" {
		if strings.TrimSpace(m.fTitle) == "" {
			return
		}
		_, err := m.svc.Add(todo.Item{
			Title: m.fTitle, Description: m.fDesc, Tags: tags, ADO: m.fADO,
			Section: m.fSection,
		})
		m.result("Added", err)
		return
	}

	err := m.svc.Edit(m.editID, func(it *todo.Item) {
		it.Title = m.fTitle
		it.Description = m.fDesc
		it.Tags = tags
		it.ADO = m.fADO
	})
	if err == nil && m.fSection != "" && m.fSection != m.editSec {
		err = m.svc.Move(m.editID, m.fSection)
	}
	m.result("Edited", err)
}

// defaultFormSection is the section a new item defaults to: the cursor's section
// if it is an open item, otherwise the first non-done section.
func (m *Model) defaultFormSection() string {
	if r, ok := m.selectedRow(); ok && r.id != "" {
		return r.section.Key
	}
	if open := m.openSections(); len(open) > 0 {
		return open[0]
	}
	return ""
}

// splitTags parses a comma-separated tag string into bare slugs.
func splitTags(s string) []string {
	var out []string
	for _, t := range strings.Split(s, ",") {
		if t = strings.TrimSpace(t); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// formWidth and formHeight size the form to the window, within sane bounds.
func formWidth(w int) int {
	switch {
	case w <= 0:
		return 60
	case w-6 < 40:
		return 40
	case w-6 > 72:
		return 72
	default:
		return w - 6
	}
}

func formHeight(h int) int {
	if h <= 0 {
		return 18
	}
	return h - 4
}
