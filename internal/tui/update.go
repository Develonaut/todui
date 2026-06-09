package tui

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
)

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if m.form != nil {
			m.form = m.form.WithWidth(formWidth(m.width)).WithHeight(formHeight(m.height))
		}
		return m, nil
	case fileChangedMsg:
		// Apply now in the list; defer while a form/confirm is open so an
		// external write never yanks data out from under an edit.
		if m.mode == modeList {
			m.rebuild()
		} else {
			m.pendingReload = true
		}
		return m, watchCmd(m.watcher, filepath.Base(m.storePath))
	case tea.KeyPressMsg:
		if m.mode == modeForm {
			return m.updateForm(msg)
		}
		return m.handleKey(msg)
	}
	if m.mode == modeForm && m.form != nil {
		return m.updateForm(msg)
	}
	return m, nil
}

// handleKey resolves a key to an action through the active scope stack.
func (m *Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	m.status = ""
	action, ok := m.keys.Match(msg.String(), m.activeScopes())
	if !ok {
		return m, nil
	}
	return m.dispatch(action)
}

// dispatch performs the behavior bound to an action via the action table.
func (m *Model) dispatch(action string) (tea.Model, tea.Cmd) {
	if h, ok := m.actions[action]; ok {
		return m, h()
	}
	return m, nil
}

// buildActions wires each action id to its behavior. Most actions mutate state
// and return no command; add/edit return the form's init command; quit returns
// tea.Quit.
func (m *Model) buildActions() map[string]func() tea.Cmd {
	none := func(fn func()) func() tea.Cmd {
		return func() tea.Cmd { fn(); return nil }
	}
	return map[string]func() tea.Cmd{
		actQuit:        func() tea.Cmd { m.closeWatcher(); return tea.Quit },
		actHelp:        none(func() { m.showHelp = !m.showHelp }),
		actReload:      none(m.rebuild),
		actUp:          none(func() { m.moveCursor(-1) }),
		actDown:        none(func() { m.moveCursor(1) }),
		actSectionPrev: none(func() { m.jumpSection(-1) }),
		actSectionNext: none(func() { m.jumpSection(1) }),
		actComplete:    none(m.complete),
		actStart:       none(m.start),
		actReorderUp:   none(func() { m.reorder(-1) }),
		actReorderDown: none(func() { m.reorder(1) }),
		actMovePrev:    none(func() { m.moveSection(-1) }),
		actMoveNext:    none(func() { m.moveSection(1) }),
		actAdd:         m.enterAdd,
		actEdit:        m.enterEdit,
		actDelete:      none(m.beginDelete),
		actConfirmYes:  none(m.confirmDelete),
		actConfirmNo:   none(m.cancelDelete),
	}
}

// updateForm drives the active Huh form and applies it on completion.
func (m *Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	fm, cmd := m.form.Update(msg)
	if f, ok := fm.(*huh.Form); ok {
		m.form = f
	}
	switch m.form.State {
	case huh.StateCompleted:
		m.applyForm()
		m.exitForm()
		return m, nil
	case huh.StateAborted:
		m.exitForm()
		return m, nil
	}
	return m, cmd
}
