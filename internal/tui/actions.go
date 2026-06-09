package tui

import "github.com/develonaut/todui/internal/todo"

// snapshot saves the current list so the change about to happen can be undone.
func (m *Model) snapshot() {
	cp := m.list
	cp.Items = append([]todo.Item(nil), m.list.Items...)
	cp.Header = append([]string(nil), m.list.Header...)
	const maxUndo = 50
	m.undo = append(m.undo, cp)
	if len(m.undo) > maxUndo {
		m.undo = m.undo[len(m.undo)-maxUndo:]
	}
}

// undoLast reverts the most recent change.
func (m *Model) undoLast() {
	if len(m.undo) == 0 {
		m.status, m.err = "nothing to undo", nil
		return
	}
	prev := m.undo[len(m.undo)-1]
	m.undo = m.undo[:len(m.undo)-1]
	m.result("Undid last change", m.svc.Replace(prev))
	m.rebuild()
}

// complete marks the selected item done.
func (m *Model) complete() {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	m.snapshot()
	m.result("Completed "+r.id, m.svc.Complete(r.id))
	m.rebuild()
}

// start claims the selected item (moves it to the start section).
func (m *Model) start() {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	m.snapshot()
	m.result("Started "+r.id, m.svc.Start(r.id))
	m.rebuild()
}

// reorder shifts the selected item within its section, and at a section
// boundary lifts it into the adjacent section (⇧↑ to the bottom of the section
// above, ⇧↓ to the top of the section below).
func (m *Model) reorder(delta int) {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	s := m.svc.Schema()
	ci := sectionIndexOf(s.Sections, r.section.Key)
	count := len(m.list.Section(s, r.section.Key))

	switch {
	case delta < 0 && r.item.Order == 0: // top of section → section above
		if ci <= 0 {
			return
		}
		dest := s.Sections[ci-1].Key
		m.snapshot()
		m.result("Moved "+r.id+" → "+s.Sections[ci-1].Title, m.svc.Move(r.id, dest))
		m.rebuild()
		m.cursorToItem(dest, r.item.Title)
	case delta > 0 && r.item.Order == count-1: // bottom of section → section below
		if ci < 0 || ci >= len(s.Sections)-1 {
			return
		}
		dest := s.Sections[ci+1].Key
		m.snapshot()
		m.result("Moved "+r.id+" → "+s.Sections[ci+1].Title, m.svc.MoveToTop(r.id, dest))
		m.rebuild()
		m.cursorToItem(dest, r.item.Title)
	default:
		m.snapshot()
		if err := m.svc.Reorder(r.id, delta); err != nil {
			m.err = err
			return
		}
		m.rebuild()
		m.moveCursor(delta)
	}
}

// moveSection moves the selected item to the adjacent section (any section,
// including out of Done).
func (m *Model) moveSection(dir int) {
	r, ok := m.currentItem()
	if !ok {
		return
	}
	secs := m.svc.Schema().Sections
	ci := -1
	for i := range secs {
		if secs[i].Key == r.section.Key {
			ci = i
			break
		}
	}
	ni := ci + dir
	if ci < 0 || ni < 0 || ni >= len(secs) {
		return
	}
	m.snapshot()
	m.result("Moved "+r.id+" → "+secs[ni].Title, m.svc.Move(r.id, secs[ni].Key))
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

// setCollapsed folds/unfolds a section and keeps the cursor on its header.
func (m *Model) setCollapsed(key string, v bool) {
	m.collapsed[key] = v
	m.rebuild()
	m.cursorToSection(key)
}

// toggleCollapse folds or unfolds the section under the cursor.
func (m *Model) toggleCollapse() {
	if r, ok := m.selectedRow(); ok {
		m.setCollapsed(r.section.Key, !m.collapsed[r.section.Key])
	}
}

// expand opens a collapsed group, or descends into an open one (tree-style →).
func (m *Model) expand() {
	r, ok := m.selectedRow()
	if !ok || !r.header {
		return
	}
	if m.collapsed[r.section.Key] {
		m.setCollapsed(r.section.Key, false)
		return
	}
	m.moveCursor(1) // already open: step into the first child
}

// collapse closes an open group, or jumps from an item to its parent group
// (tree-style ←).
func (m *Model) collapse() {
	r, ok := m.selectedRow()
	if !ok {
		return
	}
	if !r.header {
		m.cursorToSection(r.section.Key)
		return
	}
	if !m.collapsed[r.section.Key] {
		m.setCollapsed(r.section.Key, true)
	}
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
	m.snapshot()
	m.result("Deleted "+m.confirmID, m.svc.Delete(m.confirmID))
	m.confirmID = ""
	m.mode = modeList
	m.rebuild()
}
