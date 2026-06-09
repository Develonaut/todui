package todo

import (
	"testing"
	"time"
)

func TestTrimDoneByCount(t *testing.T) {
	s := testSchema()
	l := &List{}
	for i := range 5 {
		l.Items = append(l.Items, Item{Title: "d", Section: "done", DoneDate: "2026-06-0" + string(rune('1'+i))})
	}
	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if err := l.TrimDone(s, now, 3, 0); err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != 3 {
		t.Fatalf("kept %d, want 3", len(l.Items))
	}
	// newest first: 06-05, 06-04, 06-03
	if l.Items[0].DoneDate != "2026-06-05" || l.Items[2].DoneDate != "2026-06-03" {
		t.Errorf("trim order wrong: %q..%q", l.Items[0].DoneDate, l.Items[2].DoneDate)
	}
}

func TestTrimDoneByAge(t *testing.T) {
	s := testSchema()
	l := &List{Items: []Item{
		{Title: "fresh", Section: "done", DoneDate: "2026-06-08"},
		{Title: "old", Section: "done", DoneDate: "2026-05-01"},
	}}
	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if err := l.TrimDone(s, now, 10, 7); err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != 1 || l.Items[0].Title != "fresh" {
		t.Errorf("age trim kept %+v", l.Items)
	}
}

func TestTrimDoneParsesAnnotatedDate(t *testing.T) {
	s := testSchema()
	l := &List{Items: []Item{
		{Title: "annotated", Section: "done", DoneDate: "2026-06-08, /standup"},
		{Title: "old", Section: "done", DoneDate: "2026-05-01"},
	}}
	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if err := l.TrimDone(s, now, 10, 7); err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != 1 || l.Items[0].Title != "annotated" {
		t.Errorf("annotated-date trim kept %+v", l.Items)
	}
}

func TestTrimDoneLeavesOpenItems(t *testing.T) {
	s := testSchema()
	l := sample() // no done items
	now := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	before := len(l.Items)
	if err := l.TrimDone(s, now, 1, 1); err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != before {
		t.Errorf("trim removed open items: %d -> %d", before, len(l.Items))
	}
}
