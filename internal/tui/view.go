package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/develonaut/todui/internal/todo"
)

// View implements tea.Model. It composes child strings and wraps the result in
// a single full-screen tea.View.
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

// viewList renders the screen: a title bar, the scrollable sectioned list with
// truncated titles, a detail pane for the selected task, then footer and help.
func (m *Model) viewList() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	detailH := clamp(h/4, 4, 8)
	// chrome = title(1) + blank(1) + rule(1) + footer(1) + help(1)
	listH := max(h-detailH-5, 3)

	lines, cursorLine := m.listLines(w)
	m.ensureVisible(cursorLine, listH, len(lines))
	visible := padLines(window(lines, m.listOffset, listH), listH)

	rule := styleDim.Render(strings.Repeat("─", min(w, 100)))
	detail := padLines(strings.Split(m.detailPane(min(w, 100)), "\n"), detailH)

	out := []string{m.titleBar()}
	out = append(out, "")
	out = append(out, visible...)
	out = append(out, rule)
	out = append(out, detail...)
	out = append(out, m.footer(), m.helpBar(w))
	return strings.Join(out, "\n")
}

// titleBar renders the app name and last-updated stamp.
func (m *Model) titleBar() string {
	bar := styleTitle.Render("todui")
	if m.list.LastUpdated != "" {
		bar += "  " + styleDim.Render(m.list.LastUpdated)
	}
	return bar
}

// listLines builds every list line (section headers and truncated item rows)
// and reports which line the cursor is on.
func (m *Model) listLines(w int) ([]string, int) {
	s := m.svc.Schema()
	var lines []string
	cursorLine := 0
	idx := 0
	for i, sec := range s.Sections {
		if i > 0 {
			lines = append(lines, "")
		}
		items := m.list.Section(s, sec.Key)
		lines = append(lines, sectionStyle(i).Render(sec.Title)+styleDim.Render(fmt.Sprintf("  (%d)", len(items))))
		for _, it := range items {
			if idx == m.cursor {
				cursorLine = len(lines)
			}
			lines = append(lines, m.itemLine(s, sec, it, idx == m.cursor, w))
			idx++
		}
	}
	return lines, cursorLine
}

// itemLine renders one compact, single-line item: cursor, ID, claimed dot, and
// a truncated title.
func (m *Model) itemLine(s todo.Schema, sec todo.Section, it todo.Item, selected bool, w int) string {
	mark := s.ID(sec, it.Order)
	if mark == "" {
		mark = "·"
	}
	cursor := "  "
	if selected {
		cursor = styleCursor.Render("▸ ")
	}
	claimed := " "
	if it.Claimed {
		claimed = styleClaim.Render("●")
	}

	title := truncate(shortTitle(it.Task), max(10, w-9))
	if selected {
		title = styleSelected.Render(title)
	}
	return cursor + styleID.Render(fmt.Sprintf("%-3s", mark)) + " " + claimed + " " + title
}

// detailPane renders the full detail of the selected task, wrapped to width.
func (m *Model) detailPane(w int) string {
	r, ok := m.selectedRow()
	if !ok {
		return styleDim.Render("No task selected.")
	}
	it := r.item

	id := r.id
	if id == "" {
		id = "done"
	}
	bits := []string{styleID.Render(id), r.section.Title}
	if it.Claimed {
		bits = append(bits, styleClaim.Render("claimed"))
	}
	if it.DoneDate != "" {
		bits = append(bits, styleDim.Render("done "+firstDate(it.DoneDate)))
	}

	wrap := lipgloss.NewStyle().Width(w)
	parts := []string{strings.Join(bits, styleDim.Render(" · ")), wrap.Render(it.Task)}

	var meta []string
	if it.Context != "" {
		meta = append(meta, styleDim.Render("context: ")+it.Context)
	}
	if len(it.Tags) > 0 {
		var tags []string
		for _, t := range it.Tags {
			tags = append(tags, styleTag.Render("["+t+"]"))
		}
		meta = append(meta, strings.Join(tags, " "))
	}
	if it.ADO != "" {
		meta = append(meta, styleDim.Render("ref: ")+it.ADO)
	}
	if len(meta) > 0 {
		parts = append(parts, wrap.Render(strings.Join(meta, "   ")))
	}
	return strings.Join(parts, "\n")
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

// helpBar renders the context-appropriate keys from the same keymap used for
// dispatch, fitting as many as the width allows (measured ANSI-aware).
func (m *Model) helpBar(w int) string {
	var b strings.Builder
	for _, bind := range m.keys.Help(m.activeScopes()) {
		if len(bind.Keys) == 0 {
			continue
		}
		seg := styleKey.Render(keyLabel(bind.Keys[0])) + " " + styleDim.Render(bind.Help)
		sep := ""
		if b.Len() > 0 {
			sep = "  "
		}
		if lipgloss.Width(b.String()+sep+seg) > w {
			break
		}
		b.WriteString(sep + seg)
	}
	return b.String()
}

// ensureVisible adjusts the scroll offset so the cursor line stays on screen.
func (m *Model) ensureVisible(cursorLine, height, total int) {
	if cursorLine < m.listOffset {
		m.listOffset = cursorLine
	}
	if cursorLine >= m.listOffset+height {
		m.listOffset = cursorLine - height + 1
	}
	m.listOffset = clamp(m.listOffset, 0, max(0, total-height))
}
