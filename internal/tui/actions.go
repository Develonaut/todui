package tui

// complete marks the selected item done.
func (m *Model) complete() {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	m.result("Completed "+r.id, m.svc.Complete(r.id))
	m.rebuild()
}

// start claims the selected item (moves it to the start section).
func (m *Model) start() {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	m.result("Started "+r.id, m.svc.Start(r.id))
	m.rebuild()
}

// reorder shifts the selected item within its section and follows it.
func (m *Model) reorder(delta int) {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	if err := m.svc.Reorder(r.id, delta); err != nil {
		m.err = err
		return
	}
	m.rebuild()
	m.moveCursor(delta)
}

// moveSection moves the selected item to the adjacent non-done section.
func (m *Model) moveSection(dir int) {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	open := m.openSections()
	ci := -1
	for i, key := range open {
		if key == r.section.Key {
			ci = i
			break
		}
	}
	ni := ci + dir
	if ci < 0 || ni < 0 || ni >= len(open) {
		return
	}
	m.result("Moved "+r.id+" → "+open[ni], m.svc.Move(r.id, open[ni]))
	m.rebuild()
}

// openSections returns the non-done section keys in display order.
func (m *Model) openSections() []string {
	var keys []string
	for _, sec := range m.svc.Schema().Sections {
		if !sec.Done {
			keys = append(keys, sec.Key)
		}
	}
	return keys
}

// toggleCollapse folds or unfolds the section under the cursor, keeping the
// cursor on that section's header.
func (m *Model) toggleCollapse() {
	r, ok := m.selectedRow()
	if !ok {
		return
	}
	key := r.section.Key
	m.collapsed[key] = !m.collapsed[key]
	m.rebuild()
	m.cursorToSection(key)
}

// goalBy adjusts the daily goal, clamped to a sane range.
func (m *Model) goalBy(delta int) {
	m.goal = clamp(m.goal+delta, 0, 99)
}

// beginDelete enters the delete-confirmation mode for the selected item.
func (m *Model) beginDelete() {
	if r, ok := m.currentItem(); ok {
		m.confirmID, m.mode = r.id, modeConfirm
	}
}

// cancelDelete dismisses the delete confirmation, applying any reload that was
// deferred while the prompt was open.
func (m *Model) cancelDelete() {
	m.confirmID, m.mode = "", modeList
	if m.pendingReload {
		m.rebuild()
	}
}

// confirmDelete removes the item pending confirmation.
func (m *Model) confirmDelete() {
	m.result("Deleted "+m.confirmID, m.svc.Delete(m.confirmID))
	m.confirmID = ""
	m.mode = modeList
	m.rebuild()
}
