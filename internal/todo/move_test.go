package todo

import "testing"

func TestMoveToSectionAppends(t *testing.T) {
	s := testSchema()
	l := sample() // now: a,b ; next: c ; later: d
	// move "a" (idx 0) to next -> should land at bottom of next (order 1)
	if err := l.Move(s, 0, "next"); err != nil {
		t.Fatal(err)
	}
	if l.Items[0].Section != "next" || l.Items[0].Order != 1 {
		t.Errorf("after move: section=%q order=%d, want next/1", l.Items[0].Section, l.Items[0].Order)
	}
	if err := l.Move(s, 0, "nope"); err == nil {
		t.Error("Move to unknown section should error")
	}
}

func TestReorderWithinSection(t *testing.T) {
	l := sample() // now: a(0), b(1)
	// move b up
	if err := l.Reorder(1, -1); err != nil {
		t.Fatal(err)
	}
	if l.Items[0].Order != 1 || l.Items[1].Order != 0 {
		t.Errorf("orders after reorder: %d,%d want 1,0", l.Items[0].Order, l.Items[1].Order)
	}
	// reordering the top item up is a no-op (boundary)
	top := indexByTask(l, "a")
	if err := l.Reorder(top, -1); err != nil {
		t.Errorf("boundary reorder should be a no-op, got %v", err)
	}
}

func TestComplete(t *testing.T) {
	s := testSchema()
	l := sample()
	if err := l.Complete(s, 0, "2026-06-09"); err != nil {
		t.Fatal(err)
	}
	if l.Items[0].Section != "done" || l.Items[0].DoneDate != "2026-06-09" {
		t.Errorf("complete: %+v", l.Items[0])
	}
}

func indexByTask(l *List, task string) int {
	for i := range l.Items {
		if l.Items[i].Task == task {
			return i
		}
	}
	return -1
}
