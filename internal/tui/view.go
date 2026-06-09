package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/develonaut/todui/internal/todo"
)

// titleWidth caps a list title so rows stay short regardless of panel width.
const titleWidth = 64

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

// viewList renders the framed layout: a title row, a full-width TASKS panel
// stacked above a full-width DETAIL panel, and a bottom keybar.
func (m *Model) viewList() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	// title(1) + blank(1) + tasks(listH+2) + detail(detailH+2) + keybar(1) = h
	detailH := clamp(h/4, 5, 9)
	listH := max(h-detailH-7, 3)

	tasks := framePanel("TASKS", m.listBody(w-4, listH), w, styleBorderActive)
	detail := framePanel("DETAIL", m.detailBody(w-4, detailH), w, styleBorder)

	return strings.Join([]string{
		m.titleBar(w),
		"",
		tasks,
		detail,
		m.bottomBar(w),
	}, "\n")
}

// titleBar renders the app name on the left and the last-updated stamp right.
func (m *Model) titleBar(w int) string {
	left := styleTitle.Render("todui")
	right := styleDim.Render(m.list.LastUpdated)
	return spread(left, right, w)
}

// bottomBar renders contextual help (or a status/confirm message) on the left
// and the item count on the right.
func (m *Model) bottomBar(w int) string {
	right := styleDim.Render(fmt.Sprintf("%d items", len(m.rows)))
	avail := w - lipgloss.Width(right) - 2

	var left string
	switch {
	case m.mode == modeConfirm:
		left = styleConfirm.Render("delete " + m.confirmID + "? (y/n)")
	case m.err != nil:
		left = styleErr.Render(truncate(m.err.Error(), avail))
	case m.status != "":
		left = styleStatus.Render(m.status)
	default:
		left = m.helpBar(avail)
	}
	return spread(left, right, w)
}

// listBody builds the scrollable, section-grouped task lines for the TASKS panel.
func (m *Model) listBody(cw, height int) string {
	lines, cursorLine := m.listLines(cw)
	m.ensureVisible(cursorLine, height, len(lines))
	return strings.Join(padLines(window(lines, m.listOffset, height), height), "\n")
}

// listLines builds every list line (section headers and truncated item rows)
// for inner width cw and reports which line the cursor is on.
func (m *Model) listLines(cw int) ([]string, int) {
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
			lines = append(lines, m.itemLine(s, sec, it, idx == m.cursor, cw))
			idx++
		}
	}
	return lines, cursorLine
}

// itemLine renders one compact row: cursor, ID, claimed dot, and short title.
func (m *Model) itemLine(s todo.Schema, sec todo.Section, it todo.Item, selected bool, cw int) string {
	mark := s.ID(sec, it.Order)
	if mark == "" {
		mark = "·"
	}
	cursor := "  "
	if selected {
		cursor = styleCursor.Render("▸ ")
	}
	dot := " "
	if it.Claimed {
		dot = styleClaim.Render("●")
	}

	title := truncate(it.Title, max(1, min(cw-8, titleWidth)))
	style := styleItem
	if selected {
		style = styleSelect
	}
	return cursor + styleID.Render(fmt.Sprintf("%-3s", mark)) + " " + dot + " " + style.Render(title)
}

// detailBody renders the selected task's detail for the DETAIL panel: a header
// line, the title, the dimmed description, then tags/ref.
func (m *Model) detailBody(cw, height int) string {
	r, ok := m.selectedRow()
	if !ok {
		return strings.Join(padLines([]string{styleDim.Render("No task selected.")}, height), "\n")
	}
	it := r.item

	id := r.id
	if id == "" {
		id = "done"
	}
	bits := []string{styleID.Render(id), styleDim.Render(r.section.Title)}
	if it.Claimed {
		bits = append(bits, styleClaim.Render("claimed"))
	}
	if it.DoneDate != "" {
		bits = append(bits, styleDim.Render("done "+firstDate(it.DoneDate)))
	}

	lines := []string{strings.Join(bits, styleFaint.Render(" · ")), ""}
	lines = append(lines, strings.Split(styleSelect.Width(cw).Render(it.Title), "\n")...)
	if it.Description != "" {
		lines = append(lines, "")
		lines = append(lines, strings.Split(styleDetail.Width(cw).Render(it.Description), "\n")...)
	}

	var meta []string
	if len(it.Tags) > 0 {
		tags := make([]string, len(it.Tags))
		for i, t := range it.Tags {
			tags[i] = styleTag.Render(t)
		}
		meta = append(meta, strings.Join(tags, " "))
	}
	if it.ADO != "" {
		meta = append(meta, styleDim.Render("ref ")+styleDetail.Render(it.ADO))
	}
	if len(meta) > 0 {
		lines = append(lines, "")
		lines = append(lines, strings.Split(lipgloss.NewStyle().Width(cw).Render(strings.Join(meta, "   ")), "\n")...)
	}
	return strings.Join(padLines(lines, height), "\n")
}

// helpBar renders the context-appropriate keys from the same keymap used for
// dispatch, fitting as many as the width allows (measured ANSI-aware).
func (m *Model) helpBar(w int) string {
	var b strings.Builder
	for _, bind := range m.keys.Help(m.activeScopes()) {
		if len(bind.Keys) == 0 {
			continue
		}
		key := bind.HelpKey
		if key == "" {
			key = keyLabel(bind.Keys[0])
		}
		seg := styleKey.Render(key) + " " + styleDim.Render(bind.Help)
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

// spread places left and right text on one line of width w, justified apart.
func spread(left, right string, w int) string {
	gap := max(w-lipgloss.Width(left)-lipgloss.Width(right), 1)
	return left + strings.Repeat(" ", gap) + right
}
