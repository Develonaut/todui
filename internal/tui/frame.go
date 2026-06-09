package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// framePanel draws a rounded box of the given total width with a letter-spaced
// label inset into the top border (bnto style). Body lines must already fit the
// inner width (width-4); they are padded, never truncated, so styling is safe.
func framePanel(label, body string, width int, border lipgloss.Style) string {
	cw := max(width-4, 1)
	lbl := styleLabel.Render(spaceLetters(label))
	dashes := max(width-lipgloss.Width(lbl)-5, 0)

	top := border.Render("╭─ ") + lbl + border.Render(" "+strings.Repeat("─", dashes)+"╮")
	bottom := border.Render("╰" + strings.Repeat("─", width-2) + "╯")

	var b strings.Builder
	b.WriteString(top + "\n")
	for _, line := range strings.Split(body, "\n") {
		b.WriteString(border.Render("│ ") + padVisual(line, cw) + border.Render(" │") + "\n")
	}
	b.WriteString(bottom)
	return b.String()
}

// joinRows places two equal-height blocks side by side with a gap. Each block's
// lines are uniform width (panels pad their content), so a line-wise join aligns.
func joinRows(a, b string, gap int) string {
	al, bl := strings.Split(a, "\n"), strings.Split(b, "\n")
	n := max(len(al), len(bl))
	sp := strings.Repeat(" ", gap)
	out := make([]string, n)
	for i := range n {
		la, lb := "", ""
		if i < len(al) {
			la = al[i]
		}
		if i < len(bl) {
			lb = bl[i]
		}
		out[i] = la + sp + lb
	}
	return strings.Join(out, "\n")
}

// spaceLetters letter-spaces a label, e.g. "TASKS" -> "T A S K S".
func spaceLetters(s string) string {
	r := []rune(s)
	parts := make([]string, len(r))
	for i, c := range r {
		parts[i] = string(c)
	}
	return strings.Join(parts, " ")
}

// padVisual right-pads s with spaces to visual width w (ANSI-aware). It assumes
// s already fits within w.
func padVisual(s string, w int) string {
	if gap := w - lipgloss.Width(s); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}
