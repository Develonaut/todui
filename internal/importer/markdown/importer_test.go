package markdown

import (
	"os"
	"strings"
	"testing"

	render "github.com/develonaut/todui/internal/render/markdown"
	"github.com/develonaut/todui/internal/todo"
)

func legacySchema() todo.Schema {
	return todo.Schema{Sections: []todo.Section{
		{Key: "in_progress", Title: "In Progress", Letter: "P"},
		{Key: "now", Title: "Now", Letter: "N"},
		{Key: "next", Title: "Next", Letter: "X"},
		{Key: "later", Title: "Later", Letter: "L"},
		{Key: "done", Title: "Done", Done: true},
	}}
}

func fixture(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile("../../../testdata/legacy_todo.md")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// TestImportRoundTrip is the acceptance gate: importing then rendering the
// fixture must reproduce it byte-for-byte.
func TestImportRoundTrip(t *testing.T) {
	src := fixture(t)
	s := legacySchema()
	list, err := Import(src, s)
	if err != nil {
		t.Fatal(err)
	}
	got, err := render.New(s).Render(list)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(src) {
		t.Errorf("round trip mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, src)
	}
}

func TestImportExtractsFields(t *testing.T) {
	list, err := Import(fixture(t), legacySchema())
	if err != nil {
		t.Fatal(err)
	}

	if len(list.Header) != 4 || list.Header[0] != "# My Tasks" {
		t.Errorf("header = %q", list.Header)
	}
	if list.LastUpdated != "2026-06-09 09:37 (via tool)" {
		t.Errorf("last updated = %q", list.LastUpdated)
	}

	if it, ok := find(list, "Ship the thing"); !ok || !it.Claimed || strings.Join(it.Tags, ",") != "work,urgent" {
		t.Errorf("claimed item = %+v", it)
	}
	if it, ok := find(list, "wire up"); !ok || it.ADO != "#42" || strings.Join(it.Tags, ",") != "dep" {
		t.Errorf("ado item = %+v", it)
	}
	if it, ok := find(list, "Did a thing"); !ok || it.Section != "done" || it.DoneDate != "2026-06-08, via tool" {
		t.Errorf("done item = %+v", it)
	}
	if it, ok := find(list, "Plain task"); !ok || !strings.Contains(it.Title, "`[mid-text]`") {
		t.Errorf("mid-text token should stay in task: %+v", it)
	}
}

func find(list todo.List, sub string) (todo.Item, bool) {
	for _, it := range list.Items {
		if strings.Contains(it.Title, sub) {
			return it, true
		}
	}
	return todo.Item{}, false
}
