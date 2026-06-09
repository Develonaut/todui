package tui

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/develonaut/todui/internal/app"
	tomlcodec "github.com/develonaut/todui/internal/codec/toml"
	filestore "github.com/develonaut/todui/internal/store/file"
	"github.com/develonaut/todui/internal/todo"
)

type fixedClock struct{}

func (fixedClock) Now() time.Time { return time.Date(2026, 6, 9, 8, 0, 0, 0, time.UTC) }

func testModel(t *testing.T) *Model {
	t.Helper()
	dir := t.TempDir()
	storePath := filepath.Join(dir, "todo.toml")
	repo := filestore.New(filestore.Options{
		Path:  storePath,
		Codec: tomlcodec.Codec{},
	})
	schema := todo.Schema{Sections: []todo.Section{
		{Key: "now", Title: "Now", Letter: "N"},
		{Key: "next", Title: "Next", Letter: "X"},
		{Key: "done", Title: "Done", Done: true},
	}}
	svc := app.NewService(repo, fixedClock{}, app.Settings{Schema: schema, DoneMax: 10})
	if _, err := svc.Add(todo.Item{Title: "first", Section: "now"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Add(todo.Item{Title: "second", Section: "next"}); err != nil {
		t.Fatal(err)
	}
	m := New(svc, storePath, nil)
	t.Cleanup(m.closeWatcher)
	m.width, m.height = 80, 24
	return m
}

func TestViewRenders(t *testing.T) {
	m := testModel(t)
	v := m.View()
	if !v.AltScreen {
		t.Error("expected full-screen view")
	}
	for _, want := range []string{"Now", "first", "Next", "second"} {
		if !strings.Contains(v.Content, want) {
			t.Errorf("view missing %q:\n%s", want, v.Content)
		}
	}
}

func TestDispatchComplete(t *testing.T) {
	m := testModel(t)
	if !cursorTo(m, "second") {
		t.Fatal("could not place cursor on 'second'")
	}
	m.dispatch(actComplete)
	if !inDone(m, "second") {
		t.Errorf("'second' should be in done after complete:\n%+v", m.list.Items)
	}
}

func TestDispatchDeleteWithConfirm(t *testing.T) {
	m := testModel(t)
	if !cursorTo(m, "first") {
		t.Fatal("could not place cursor on 'first'")
	}
	m.dispatch(actDelete)
	if m.mode != modeConfirm || m.confirmID == "" {
		t.Fatalf("delete should enter confirm mode")
	}
	m.dispatch(actConfirmYes)
	if m.mode != modeList {
		t.Errorf("confirm yes should return to list")
	}
	for _, it := range m.list.Items {
		if it.Title == "first" {
			t.Errorf("'first' should be deleted")
		}
	}
}

func TestCollapseFoldsSection(t *testing.T) {
	m := testModel(t)
	m.cursor = 0 // first row is the "now" header
	if got := m.activeScopes(); got[0] != scopeHeader {
		t.Fatalf("on a header, top scope = %q, want header", got[0])
	}
	before := len(m.rows)
	m.dispatch(actCollapse)
	if !m.collapsed["now"] || len(m.rows) >= before {
		t.Errorf("collapsing 'now' should hide its items (rows %d -> %d)", before, len(m.rows))
	}
}

func TestActiveScopesByContext(t *testing.T) {
	m := testModel(t)
	if got := m.activeScopes(); got[0] != scopeItem { // cursor starts on first item
		t.Errorf("on an item, top scope = %q, want item", got[0])
	}
	m.rows = nil
	if got := m.activeScopes(); got[0] != scopeEmpty {
		t.Errorf("with no rows, top scope = %q, want empty", got[0])
	}
}

func cursorTo(m *Model, title string) bool {
	for i, r := range m.rows {
		if !r.header && r.item.Title == title {
			m.cursor = i
			return true
		}
	}
	return false
}

func inDone(m *Model, title string) bool {
	for _, it := range m.list.Items {
		if it.Title == title && it.Section == "done" {
			return true
		}
	}
	return false
}
