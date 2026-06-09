package tui

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
