package tui

import (
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/fsnotify/fsnotify"
)

// fileChangedMsg signals that the store file changed on disk (e.g. an external
// CLI invocation or the /todo skill wrote it).
type fileChangedMsg struct{}

// newWatcher watches the store file's directory (watching the directory, not the
// file, survives the atomic rename that replaces the inode). It returns nil,
// disabling live reload, if a watcher cannot be established.
func newWatcher(storePath string) *fsnotify.Watcher {
	dir := filepath.Dir(storePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil
	}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}
	if err := w.Add(dir); err != nil {
		_ = w.Close()
		return nil
	}
	return w
}

// watchCmd blocks until the store file changes, then emits fileChangedMsg. It is
// re-issued after each event to keep listening.
func watchCmd(w *fsnotify.Watcher, base string) tea.Cmd {
	if w == nil {
		return nil
	}
	return func() tea.Msg {
		for ev := range w.Events {
			if filepath.Base(ev.Name) == base && ev.Has(fsnotify.Write|fsnotify.Create|fsnotify.Rename) {
				return fileChangedMsg{}
			}
		}
		return nil
	}
}

// closeWatcher stops the watcher; it is safe to call more than once.
func (m *Model) closeWatcher() {
	if m.watcher != nil {
		_ = m.watcher.Close()
		m.watcher = nil
	}
}
