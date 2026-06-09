// Package tui is the interactive Bubble Tea front-end. It is a thin driver over
// the application service: it renders the list, dispatches keys through the
// keymap, and runs Huh forms for add/edit. All mutations go through app.Service.
package tui

import (
	"path/filepath"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
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

// visRow is one navigable line: either a section header or an item. Headers are
// selectable so they can be collapsed/expanded.
type visRow struct {
	header  bool
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
	listOffset int             // first visible list line (vertical scroll)
	collapsed  map[string]bool // section key -> collapsed

	goal     int            // daily completion goal
	progress progress.Model // header progress bar

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
		collapsed: map[string]bool{},
		goal:      svc.Goal(),
		progress:  progress.New(),
	}
	if dk := svc.Schema().DoneKey(); dk != "" {
		m.collapsed[dk] = true // start with completed items folded away
	}
	m.actions = m.buildActions()
	m.rebuild()
	m.cursor = m.firstItemIndex()
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

// rebuild reloads the list from the service and flattens it into navigable rows
// (a header per section, then its items unless the section is collapsed).
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
		m.rows = append(m.rows, visRow{header: true, section: sec, secIdx: i})
		if m.collapsed[sec.Key] {
			continue
		}
		for _, it := range l.Section(s, sec.Key) {
			m.rows = append(m.rows, visRow{item: it, id: s.ID(sec, it.Order), section: sec, secIdx: i})
		}
	}
	m.cursor = clamp(m.cursor, 0, max(0, len(m.rows)-1))
	m.pendingReload = false
}

// firstItemIndex returns the row of the first item (skipping leading headers).
func (m *Model) firstItemIndex() int {
	for i := range m.rows {
		if !m.rows[i].header {
			return i
		}
	}
	return 0
}

// selectedRow returns the row under the cursor, if any.
func (m *Model) selectedRow() (visRow, bool) {
	if m.cursor < 0 || m.cursor >= len(m.rows) {
		return visRow{}, false
	}
	return m.rows[m.cursor], true
}

// currentItem returns the selected row only when it is an addressable item.
func (m *Model) currentItem() (visRow, bool) {
	r, ok := m.selectedRow()
	if !ok || r.header || r.id == "" {
		return visRow{}, false
	}
	return r, true
}

// cursorToSection moves the cursor onto a section's header.
func (m *Model) cursorToSection(key string) {
	for i := range m.rows {
		if m.rows[i].header && m.rows[i].section.Key == key {
			m.cursor = i
			return
		}
	}
}

// doneToday counts items completed on today's date.
func (m *Model) doneToday() int {
	today := time.Now().Format("2006-01-02")
	n := 0
	for i := range m.list.Items {
		if strings.HasPrefix(m.list.Items[i].DoneDate, today) {
			n++
		}
	}
	return n
}

// result records the outcome of a mutation as a status message or an error.
func (m *Model) result(msg string, err error) {
	if err != nil {
		m.err = err
		return
	}
	m.status, m.err = msg, nil
}
