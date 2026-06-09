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
	repo := filestore.New(filestore.Options{
		Path:  filepath.Join(dir, "todo.toml"),
		Codec: tomlcodec.Codec{},
	})
	schema := todo.Schema{Sections: []todo.Section{
		{Key: "now", Title: "Now", Letter: "N"},
		{Key: "next", Title: "Next", Letter: "X"},
		{Key: "done", Title: "Done", Done: true},
	}}
	svc := app.NewService(repo, fixedClock{}, app.Settings{Schema: schema, DoneMax: 10})
	if _, err := svc.Add(todo.Item{Task: "first", Section: "now"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Add(todo.Item{Task: "second", Section: "next"}); err != nil {
		t.Fatal(err)
	}
	m := New(svc, nil)
	m.width, m.height = 80, 24
	return m
}

func TestViewRenders(t *testing.T) {
	m := testModel(t)
	v := m.View()
	if !v.AltScreen {
		t.Error("expected full-screen view")
	}
	for _, want := range []string{"todui", "Now", "first", "Next", "second"} {
		if !strings.Contains(v.Content, want) {
			t.Errorf("view missing %q:\n%s", want, v.Content)
		}
	}
}

func TestDispatchNavigateAndComplete(t *testing.T) {
	m := testModel(t)
	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0", m.cursor)
	}
	m.dispatch(actDown)
	if m.cursor != 1 || m.rows[m.cursor].item.Task != "second" {
		t.Fatalf("after down: cursor=%d task=%q", m.cursor, m.rows[m.cursor].item.Task)
	}
	m.dispatch(actComplete)
	if it, ok := find(m, "second"); !ok || it.section.Key != "done" {
		t.Errorf("second should be done after complete")
	}
}

func TestDispatchDeleteWithConfirm(t *testing.T) {
	m := testModel(t)
	m.dispatch(actDelete)
	if m.mode != modeConfirm || m.confirmID == "" {
		t.Fatalf("delete should enter confirm mode")
	}
	m.dispatch(actConfirmYes)
	if m.mode != modeList {
		t.Errorf("confirm yes should return to list")
	}
	if _, ok := find(m, "first"); ok {
		t.Errorf("first should be deleted")
	}
}

func TestActiveScopesByContext(t *testing.T) {
	m := testModel(t)
	if got := m.activeScopes(); got[0] != scopeItem {
		t.Errorf("with items, top scope = %q, want item", got[0])
	}
	m.rows = nil
	if got := m.activeScopes(); got[0] != scopeEmpty {
		t.Errorf("with no items, top scope = %q, want empty", got[0])
	}
}

func find(m *Model, sub string) (visRow, bool) {
	for _, r := range m.rows {
		if strings.Contains(r.item.Task, sub) {
			return r, true
		}
	}
	return visRow{}, false
}
