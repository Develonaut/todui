package app

import (
	"testing"
	"time"

	"github.com/develonaut/todui/internal/todo"
)

// memRepo is an in-memory ports.Repository for tests.
type memRepo struct{ list todo.List }

func (m *memRepo) Load() (todo.List, error) { return m.list, nil }

func (m *memRepo) Update(fn func(*todo.List) error) error {
	l := m.list
	l.Items = append([]todo.Item(nil), m.list.Items...)
	if err := fn(&l); err != nil {
		return err
	}
	m.list = l
	return nil
}

type fakeClock struct{ t time.Time }

func (f fakeClock) Now() time.Time { return f.t }

func testSvc(start string) (*Service, *memRepo) {
	repo := &memRepo{}
	set := Settings{
		Schema: todo.Schema{Sections: []todo.Section{
			{Key: "now", Title: "Now", Letter: "N"},
			{Key: "next", Title: "Next", Letter: "X"},
			{Key: "in_progress", Title: "In Progress", Letter: "P"},
			{Key: "done", Title: "Done", Done: true},
		}},
		StartSection: start,
		DoneMax:      10,
	}
	clk := fakeClock{t: time.Date(2026, 6, 9, 8, 30, 0, 0, time.UTC)}
	return NewService(repo, clk, set), repo
}

func TestAddReturnsSequentialIDs(t *testing.T) {
	svc, _ := testSvc("")
	id1, err := svc.Add(todo.Item{Title: "a", Section: "now"})
	if err != nil || id1 != "NA" {
		t.Fatalf("first add = %q,%v want NA", id1, err)
	}
	id2, _ := svc.Add(todo.Item{Title: "b", Section: "now"})
	if id2 != "NB" {
		t.Errorf("second add = %q want NB", id2)
	}
}

func TestCompleteStampsDateAndTime(t *testing.T) {
	svc, _ := testSvc("")
	if _, err := svc.Add(todo.Item{Title: "a", Section: "now"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Complete("NA"); err != nil {
		t.Fatal(err)
	}
	l, _ := svc.List()
	if len(l.Items) != 1 || l.Items[0].Section != "done" || l.Items[0].DoneDate != "2026-06-09" {
		t.Errorf("after complete: %+v", l.Items)
	}
	if l.LastUpdated != "2026-06-09 08:30" {
		t.Errorf("LastUpdated = %q", l.LastUpdated)
	}
}

func TestStartWithoutStartSectionJustClaims(t *testing.T) {
	svc, _ := testSvc("")
	if _, err := svc.Add(todo.Item{Title: "a", Section: "now"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Start("NA"); err != nil {
		t.Fatal(err)
	}
	l, _ := svc.List()
	if l.Items[0].Section != "now" || !l.Items[0].Claimed {
		t.Errorf("start (no start section): %+v", l.Items[0])
	}
}

func TestStartMovesToStartSection(t *testing.T) {
	svc, _ := testSvc("in_progress")
	if _, err := svc.Add(todo.Item{Title: "a", Section: "now"}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Start("NA"); err != nil {
		t.Fatal(err)
	}
	l, _ := svc.List()
	if l.Items[0].Section != "in_progress" || !l.Items[0].Claimed {
		t.Errorf("start (with start section): %+v", l.Items[0])
	}
}

func TestUnknownIDErrors(t *testing.T) {
	svc, _ := testSvc("")
	if err := svc.Complete("ZZ"); err == nil {
		t.Error("completing unknown id should error")
	}
}
