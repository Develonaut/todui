package file

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	tomlcodec "github.com/develonaut/todui/internal/codec/toml"
	mdrender "github.com/develonaut/todui/internal/render/markdown"
	"github.com/develonaut/todui/internal/todo"
)

func testSchema() todo.Schema {
	return todo.Schema{Sections: []todo.Section{
		{Key: "now", Title: "Now", Letter: "N"},
		{Key: "done", Title: "Done", Done: true},
	}}
}

func newRepo(t *testing.T, withMirror bool) *Repository {
	t.Helper()
	dir := t.TempDir()
	opts := Options{Path: filepath.Join(dir, "todo.toml"), Codec: tomlcodec.Codec{}}
	if withMirror {
		opts.Renderer = mdrender.New(testSchema())
		opts.Mirror = filepath.Join(dir, "TODO.md")
	}
	return New(opts)
}

func TestMissingFileIsEmpty(t *testing.T) {
	l, err := newRepo(t, false).Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != 0 {
		t.Errorf("want empty list, got %d items", len(l.Items))
	}
}

func TestUpdateAndLoad(t *testing.T) {
	r := newRepo(t, false)
	err := r.Update(func(l *todo.List) error {
		l.Add(todo.Item{Task: "x", Section: "now"})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	l, err := r.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != 1 || l.Items[0].Task != "x" {
		t.Errorf("unexpected items: %+v", l.Items)
	}
}

func TestMirrorWrittenReadOnly(t *testing.T) {
	r := newRepo(t, true)
	err := r.Update(func(l *todo.List) error {
		l.Add(todo.Item{Task: "x", Section: "now"})
		l.Normalize(testSchema())
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(r.opts.Mirror)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o444 {
		t.Errorf("mirror mode = %v, want 0444", info.Mode().Perm())
	}
	b, err := os.ReadFile(r.opts.Mirror)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "**NA** x") {
		t.Errorf("mirror missing rendered item:\n%s", b)
	}
}

func TestConcurrentUpdatesDoNotLoseWrites(t *testing.T) {
	r := newRepo(t, false)
	const n = 20
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = r.Update(func(l *todo.List) error {
				l.Add(todo.Item{Task: "t", Section: "now"})
				return nil
			})
		}()
	}
	wg.Wait()

	l, err := r.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != n {
		t.Errorf("got %d items, want %d (lost updates)", len(l.Items), n)
	}
}
