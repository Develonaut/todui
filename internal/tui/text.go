package tui

import "strings"

// Text and layout helpers for the view. These operate on plain strings; callers
// apply styling after truncation so ANSI escapes are never cut mid-sequence.

// clamp constrains v to [lo, hi].
func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// window returns up to h lines starting at off.
func window(lines []string, off, h int) []string {
	off = clamp(off, 0, len(lines))
	end := min(off+h, len(lines))
	return lines[off:end]
}

// padLines returns exactly n lines, truncating or padding with blanks.
func padLines(lines []string, n int) []string {
	if len(lines) > n {
		return lines[:n]
	}
	for len(lines) < n {
		lines = append(lines, "")
	}
	return lines
}

// truncate shortens a plain string to at most n runes, adding an ellipsis.
func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	if n == 1 {
		return "…"
	}
	return string(r[:n-1]) + "…"
}

// shortTitle returns the leading clause of a task — the text before the first
// " — " separator — as a compact list title.
func shortTitle(task string) string {
	if i := strings.Index(task, " — "); i > 0 {
		return task[:i]
	}
	return task
}

// firstDate returns the leading YYYY-MM-DD of a done annotation.
func firstDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
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
