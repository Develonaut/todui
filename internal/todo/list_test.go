package todo

import "testing"

func TestNormalizeSortsAndRenumbers(t *testing.T) {
	s := testSchema()
	l := &List{Items: []Item{
		{Title: "later1", Section: "later", Order: 5},
		{Title: "now2", Section: "now", Order: 9},
		{Title: "now1", Section: "now", Order: 2},
		{Title: "ip", Section: "in_progress", Order: 0},
	}}
	l.Normalize(s)
	// in_progress comes before now before later; orders dense from 0.
	if l.Items[0].Title != "ip" {
		t.Errorf("first item = %q, want ip", l.Items[0].Title)
	}
	now := l.Section(s, "now")
	if len(now) != 2 || now[0].Title != "now1" || now[0].Order != 0 || now[1].Order != 1 {
		t.Errorf("now section not normalized: %+v", now)
	}
}
