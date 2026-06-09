// Package tui is the interactive Bubble Tea front-end. It is a thin driver over
// the application service: it renders the list, dispatches keys through the
// keymap, and runs Huh forms for add/edit. All mutations go through app.Service.
package tui

import (
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	huh "charm.land/huh/v2"
	"github.com/fsnotify/fsnotify"

	"github.com/develonaut/todui/internal/app"
	"github.com/develonaut/todui/internal/keymap"
	"github.com/develonaut/todui/internal/todo"
)

type mode int

const (
	modeList mode = iota
	modeForm
	modeConfirm
)

// visRow is one rendered item plus the context needed to act on it.
type visRow struct {
	item    todo.Item
	id      string
	section todo.Section
	secIdx  int
}

// Model is the Bubble Tea model.
type Model struct {
	svc  *app.Service
	keys *keymap.Keymap

	width, height int

	list       todo.List
	rows       []visRow
	cursor     int
	listOffset int // first visible list line (vertical scroll)

	mode      mode
	form      *huh.Form
	editID    string
	editSec   string
	confirmID string
	showHelp  bool

	// actions maps an action id to its behavior (paired with the keymap, which
	// maps keys to action ids).
	actions map[string]func() tea.Cmd

	// live reload
	storePath     string
	watcher       *fsnotify.Watcher
	pendingReload bool

	// form field bindings
	fTitle, fDesc, fTags, fADO, fSection string
	fClaimed                             bool

	status string
	err    error
}

// New builds the model, applying any keybinding overrides over the defaults and
// establishing a watcher on the store for live reload.
func New(svc *app.Service, storePath string, overrides []keymap.Override) *Model {
	km := defaultKeymap()
	if len(overrides) > 0 {
		km, _ = km.Merge(overrides)
	}
	m := &Model{
		svc:       svc,
		keys:      km,
		mode:      modeList,
		storePath: storePath,
		watcher:   newWatcher(storePath),
	}
	m.actions = m.buildActions()
	m.rebuild()
	return m
}

// Run launches the interactive program.
func Run(svc *app.Service, storePath string, overrides []keymap.Override) error {
	m := New(svc, storePath, overrides)
	defer m.closeWatcher()
	_, err := tea.NewProgram(m).Run()
	return err
}

// Init implements tea.Model: it begins listening for external file changes.
func (m *Model) Init() tea.Cmd {
	return watchCmd(m.watcher, filepath.Base(m.storePath))
}

// rebuild reloads the list from the service and flattens it into rows.
func (m *Model) rebuild() {
	l, err := m.svc.List()
	if err != nil {
		m.err = err
		return
	}
	m.list = l
	s := m.svc.Schema()
	m.rows = m.rows[:0]
	for i, sec := range s.Sections {
		for _, it := range l.Section(s, sec.Key) {
			m.rows = append(m.rows, visRow{item: it, id: s.ID(sec, it.Order), section: sec, secIdx: i})
		}
	}
	if m.cursor >= len(m.rows) {
		m.cursor = len(m.rows) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.pendingReload = false
}

// selectedRow returns the row under the cursor, if any.
func (m *Model) selectedRow() (visRow, bool) {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return visRow{}, false
	}
	return m.rows[m.cursor], true
}

// result records the outcome of a mutation as a status message or an error.
func (m *Model) result(msg string, err error) {
	if err != nil {
		m.err = err
		return
	}
	m.status, m.err = msg, nil
}
