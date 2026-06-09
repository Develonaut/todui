package tui

import "github.com/develonaut/todui/internal/todo"

// sectionIndexOf returns the display index of a section, or -1.
func sectionIndexOf(secs []todo.Section, key string) int {
	for i := range secs {
		if secs[i].Key == key {
			return i
		}
	}
	return -1
}

// cursorToItem places the cursor on an item matched by section and title.
func (m *Model) cursorToItem(section, title string) {
	for i := range m.rows {
		if !m.rows[i].header && m.rows[i].section.Key == section && m.rows[i].item.Title == title {
			m.cursor = i
			return
		}
	}
}

// moveCursor shifts the cursor by delta within the row list, clamped.
func (m *Model) moveCursor(delta int) {
	m.cursor = max(0, min(m.cursor+delta, len(m.rows)-1))
}

// jumpSection moves the cursor to the first row of the next (dir>0) or previous
// (dir<0) populated section.
func (m *Model) jumpSection(dir int) {
	if len(m.rows) == 0 {
		return
	}
	starts := m.sectionStarts()
	cur := m.rows[m.cursor].section.Key
	ci := 0
	for i, st := range starts {
		if st.key == cur {
			ci = i
			break
		}
	}
	ni := ci + dir
	if ni < 0 || ni >= len(starts) {
		return
	}
	m.cursor = starts[ni].first
}

type sectionStart struct {
	key   string
	first int
}

// sectionStarts lists each populated section and the row index it begins at.
func (m *Model) sectionStarts() []sectionStart {
	var starts []sectionStart
	last := ""
	for i, r := range m.rows {
		if r.section.Key != last {
			starts = append(starts, sectionStart{r.section.Key, i})
			last = r.section.Key
		}
	}
	return starts
}
