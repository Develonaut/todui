package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	header := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FD2")).Bold(true).Render(title)
	body := lipgloss.JoinVertical(lipgloss.Left, header, "", m.form.View())
	return lipgloss.NewStyle().Margin(1, 2).Render(body)
}

// viewList renders the logo + goal bar, a full-width TASKS panel above a DETAIL
// panel, and a bottom keybar.
func (m *Model) viewList() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	// logo(1) + goal(1) + blank(1) + tasks(listH+2) + detail(detailH+2) + keybar(1)
	detailH := clamp(h/4, 5, 9)
	listH := max(h-detailH-8, 3)

	tasks := framePanel("TASKS", m.listBody(w-4, listH), w, styleBorderActive)
	detail := framePanel("DETAIL", m.detailBody(w-4, detailH), w, styleBorder)

	return strings.Join([]string{
		spread(logo(), styleDim.Render(m.list.LastUpdated), w),
		m.goalBar(w),
		"",
		tasks,
		detail,
		m.bottomBar(w),
	}, "\n")
}

// logo renders the TODUI wordmark with a magenta→purple per-letter gradient.
func logo() string {
	const letters = "TODUI"
	var b strings.Builder
	for i, ch := range letters {
		c := logoColors[i%len(logoColors)]
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Bold(true).Render(string(ch)))
		if i < len(letters)-1 {
			b.WriteByte(' ')
		}
	}
	return b.String()
}

// goalBar renders the daily-goal progress bar and count (empty when no goal).
func (m *Model) goalBar(w int) string {
	if m.goal <= 0 {
		return ""
	}
	done := m.doneToday()
	pct := float64(done) / float64(m.goal)
	if pct > 1 {
		pct = 1
	}
	return m.progress.ViewAs(pct) + styleDim.Render(fmt.Sprintf("  %d/%d done today", done, m.goal))
}

// bottomBar renders contextual help (or a status/confirm message) on the left
// and the item count on the right.
func (m *Model) bottomBar(w int) string {
	right := styleDim.Render(fmt.Sprintf("%d items", len(m.list.Items)))
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

// listBody builds the scrollable, section-grouped lines for the TASKS panel.
func (m *Model) listBody(cw, height int) string {
	lines, cursorLine := m.listLines(cw)
	m.ensureVisible(cursorLine, height, len(lines))
	return strings.Join(padLines(window(lines, m.listOffset, height), height), "\n")
}

// listLines renders every navigable row (headers and items) and reports which
// line the cursor is on.
func (m *Model) listLines(cw int) ([]string, int) {
	var lines []string
	cursorLine := 0
	for i, r := range m.rows {
		if r.header {
			if i > 0 {
				lines = append(lines, "")
			}
			if i == m.cursor {
				cursorLine = len(lines)
			}
			lines = append(lines, m.headerLine(r, i == m.cursor))
			continue
		}
		if i == m.cursor {
			cursorLine = len(lines)
		}
		lines = append(lines, m.itemLine(r, cw, i == m.cursor))
	}
	return lines, cursorLine
}

// headerLine renders a section header with a fold indicator and count.
func (m *Model) headerLine(r visRow, selected bool) string {
	icon := "▾"
	if m.collapsed[r.section.Key] {
		icon = "▸"
	}
	count := len(m.list.Section(m.svc.Schema(), r.section.Key))
	iconStyle, titleStyle := styleDim, sectionStyle(r.secIdx)
	if selected {
		iconStyle, titleStyle = styleCursor, titleStyle.Bold(true)
	}
	return iconStyle.Render(icon) + " " + titleStyle.Render(r.section.Title) + styleDim.Render(fmt.Sprintf("  (%d)", count))
}

// itemLine renders one compact row: indent, cursor, ID, and short title.
func (m *Model) itemLine(r visRow, cw int, selected bool) string {
	mark := r.id
	if mark == "" {
		mark = "·"
	}
	cur := "  "
	if selected {
		cur = styleCursor.Render("▸ ")
	}
	title := truncate(r.item.Title, max(1, min(cw-9, titleWidth)))
	style := styleItem
	if selected {
		style = styleSelect
	}
	return "  " + cur + styleID.Render(fmt.Sprintf("%-3s", mark)) + " " + style.Render(title)
}

// detailBody renders the selected row's detail: a section summary for a header,
// or the item's header line, title, dimmed description, and tags/ref.
func (m *Model) detailBody(cw, height int) string {
	r, ok := m.selectedRow()
	if !ok {
		return strings.Join(padLines([]string{styleDim.Render("Nothing selected.")}, height), "\n")
	}

	if r.header {
		count := len(m.list.Section(m.svc.Schema(), r.section.Key))
		state := "expanded"
		if m.collapsed[r.section.Key] {
			state = "collapsed"
		}
		lines := []string{
			sectionStyle(r.secIdx).Bold(true).Render(r.section.Title),
			"",
			styleDim.Render(fmt.Sprintf("%d item(s) · %s · space to fold", count, state)),
		}
		return strings.Join(padLines(lines, height), "\n")
	}

	it := r.item
	id := r.id
	if id == "" {
		id = "·"
	}
	bits := []string{styleID.Render(id), styleDim.Render(r.section.Title)}
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
