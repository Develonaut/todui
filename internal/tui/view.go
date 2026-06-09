package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/develonaut/todui/internal/todo"
)

// View implements tea.Model. It composes child strings and wraps the result in
// a single tea.View (full-screen).
func (m *Model) View() tea.View {
	var content string
	if m.mode == modeForm && m.form != nil {
		content = m.viewForm()
	} else {
		content = m.viewList()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// viewForm renders the active add/edit form.
func (m *Model) viewForm() string {
	title := "New task"
	if m.editID != "" {
		title = "Edit " + m.editID
	}
	body := lipgloss.JoinVertical(lipgloss.Left, styleTitle.Render(title), "", m.form.View())
	return lipgloss.NewStyle().Margin(1, 2).Render(body)
}

// viewList renders the sectioned task list, footer, and help bar.
func (m *Model) viewList() string {
	s := m.svc.Schema()
	var b strings.Builder

	b.WriteString(styleTitle.Render("todui"))
	if m.list.LastUpdated != "" {
		b.WriteString("  " + styleDim.Render(m.list.LastUpdated))
	}
	b.WriteString("\n\n")

	idx := 0
	for i, sec := range s.Sections {
		items := m.list.Section(s, sec.Key)
		b.WriteString(sectionStyle(i).Render(sec.Title))
		b.WriteString(styleDim.Render(fmt.Sprintf("  (%d)", len(items))))
		b.WriteByte('\n')
		for _, it := range items {
			b.WriteString(m.renderRow(s, sec, it, idx == m.cursor))
			b.WriteByte('\n')
			idx++
		}
		b.WriteByte('\n')
	}

	b.WriteString(m.footer())
	b.WriteByte('\n')
	b.WriteString(m.helpBar())
	return b.String()
}

// renderRow renders one item line.
func (m *Model) renderRow(s todo.Schema, sec todo.Section, it todo.Item, selected bool) string {
	mark := s.ID(sec, it.Order)
	if mark == "" {
		mark = "·"
	}
	cursor := "  "
	if selected {
		cursor = styleCursor.Render("▸ ")
	}
	idCol := styleID.Render(fmt.Sprintf("%-3s", mark))

	body := it.Task
	if it.Claimed {
		body += " " + styleClaim.Render("CLAIMED")
	}
	for _, t := range it.Tags {
		body += " " + styleTag.Render("["+t+"]")
	}
	if it.DoneDate != "" {
		body += " " + styleDim.Render("(done "+firstDate(it.DoneDate)+")")
	}
	if selected {
		body = styleSelected.Render(body)
	}
	return cursor + idCol + " " + body
}

// footer renders the transient status / confirm / error line.
func (m *Model) footer() string {
	switch {
	case m.mode == modeConfirm:
		return styleConfirm.Render(fmt.Sprintf("delete %s? (y/n)", m.confirmID))
	case m.err != nil:
		return styleErr.Render("error: " + m.err.Error())
	case m.status != "":
		return styleStatus.Render(m.status)
	default:
		return styleDim.Render(fmt.Sprintf("%d items", len(m.rows)))
	}
}

// helpBar renders the context-appropriate keys, generated from the same keymap
// used for dispatch.
func (m *Model) helpBar() string {
	var parts []string
	for _, b := range m.keys.Help(m.activeScopes()) {
		if len(b.Keys) == 0 {
			continue
		}
		parts = append(parts, styleKey.Render(keyLabel(b.Keys[0]))+" "+styleDim.Render(b.Help))
	}
	return strings.Join(parts, "  ")
}

// keyLabel prettifies a key string for display.
func keyLabel(k string) string {
	switch k {
	case " ", "space":
		return "space"
	case "up":
		return "↑"
	case "down":
		return "↓"
	case "left":
		return "←"
	case "right":
		return "→"
	default:
		return k
	}
}

// firstDate returns the leading YYYY-MM-DD of a done annotation.
func firstDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}
