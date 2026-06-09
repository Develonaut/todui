package markdown

import (
	"strings"
	"testing"

	"github.com/develonaut/todui/internal/todo"
)

func schema() todo.Schema {
	return todo.Schema{Sections: []todo.Section{
		{Key: "in_progress", Title: "In Progress", Letter: "P"},
		{Key: "now", Title: "Now", Letter: "N"},
		{Key: "next", Title: "Next", Letter: "X"},
		{Key: "later", Title: "Later", Letter: "L"},
		{Key: "done", Title: "Done", Done: true},
	}}
}

func TestRender(t *testing.T) {
	s := schema()
	l := todo.List{
		Header:      []string{"# TODO", "", "_subtitle_"},
		LastUpdated: "2026-06-09 08:50",
		Items: []todo.Item{
			{Task: "claimed one", Section: "in_progress", Claimed: true},
			{Task: "first", Section: "now", ADO: "#1", Tags: []string{"x"}},
			{Task: "second", Section: "now"},
			{Task: "ctx", Section: "next", Context: "because"},
			{Task: "shipped", Section: "done", DoneDate: "2026-06-04"},
		},
	}
	l.Normalize(s)

	want := strings.Join([]string{
		"# TODO",
		"",
		"_subtitle_",
		"_Last updated: 2026-06-09 08:50_",
		"",
		"## In Progress",
		"- [ ] **PA** **CLAIMED** claimed one",
		"",
		"## Now",
		"- [ ] **NA** #1 first `[x]`",
		"- [ ] **NB** second",
		"",
		"## Next",
		"- [ ] **XA** ctx — because",
		"",
		"## Later",
		"",
		"## Done",
		"- [x] shipped (done 2026-06-04)",
	}, "\n") + "\n"

	got, err := New(s).Render(l)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Errorf("render mismatch:\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
