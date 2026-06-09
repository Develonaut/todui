package todo

import "testing"

func TestAddAppendsToSection(t *testing.T) {
	l := sample()
	idx := l.Add(Item{Title: "e", Section: "now"})
	if l.Items[idx].Title != "e" {
		t.Fatalf("Add returned idx %d for wrong item", idx)
	}
	if l.Items[idx].Order != 2 {
		t.Errorf("new now item Order = %d, want 2 (appended)", l.Items[idx].Order)
	}
}

func TestEdit(t *testing.T) {
	l := sample()
	if err := l.Edit(0, func(it *Item) { it.Title = "edited"; it.Tags = []string{"x"} }); err != nil {
		t.Fatal(err)
	}
	if l.Items[0].Title != "edited" || len(l.Items[0].Tags) != 1 {
		t.Errorf("edit not applied: %+v", l.Items[0])
	}
	if err := l.Edit(99, func(*Item) {}); err == nil {
		t.Error("Edit out of range should error")
	}
}

func TestDelete(t *testing.T) {
	l := sample()
	n := len(l.Items)
	if err := l.Delete(0); err != nil {
		t.Fatal(err)
	}
	if len(l.Items) != n-1 {
		t.Errorf("len = %d, want %d", len(l.Items), n-1)
	}
	if err := l.Delete(99); err == nil {
		t.Error("Delete out of range should error")
	}
}
