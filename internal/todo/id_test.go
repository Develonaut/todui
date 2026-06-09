package todo

import "testing"

func TestSeqLabelRoundTrip(t *testing.T) {
	cases := []struct {
		n     int
		label string
	}{
		{0, "A"}, {1, "B"}, {25, "Z"}, {26, "AA"}, {27, "AB"}, {51, "AZ"}, {52, "BA"},
	}
	for _, c := range cases {
		if got := seqLabel(c.n); got != c.label {
			t.Errorf("seqLabel(%d) = %q, want %q", c.n, got, c.label)
		}
		got, ok := seqIndex(c.label)
		if !ok || got != c.n {
			t.Errorf("seqIndex(%q) = %d,%v, want %d,true", c.label, got, ok, c.n)
		}
	}
}

func TestSeqIndexInvalid(t *testing.T) {
	for _, s := range []string{"", "a", "1", "A1", " A"} {
		if _, ok := seqIndex(s); ok {
			t.Errorf("seqIndex(%q) should be invalid", s)
		}
	}
}

func TestComputeID(t *testing.T) {
	s := testSchema()
	l := &List{Items: []Item{
		{Task: "a", Section: "now"},
		{Task: "b", Section: "now"},
		{Task: "c", Section: "next"},
		{Task: "done one", Section: "done", DoneDate: "2026-06-01"},
	}}
	l.Normalize(s)
	want := map[string]string{"a": "NA", "b": "NB", "c": "XA", "done one": ""}
	for i := range l.Items {
		got := l.ComputeID(s, i)
		if got != want[l.Items[i].Task] {
			t.Errorf("ComputeID(%q) = %q, want %q", l.Items[i].Task, got, want[l.Items[i].Task])
		}
	}
}

func TestResolveRoundTrip(t *testing.T) {
	s := testSchema()
	l := sample()
	for i := range l.Items {
		id := l.ComputeID(s, i)
		if id == "" {
			continue
		}
		// IDs resolve case-insensitively and with surrounding space.
		idx, err := l.Resolve(s, "  "+id+"  ")
		if err != nil || idx != i {
			t.Errorf("Resolve(%q) = %d,%v, want %d", id, idx, err, i)
		}
	}
}

func TestResolveUnknown(t *testing.T) {
	s := testSchema()
	l := sample()
	for _, id := range []string{"ZZ", "N9", "NZ", ""} {
		if _, err := l.Resolve(s, id); err == nil {
			t.Errorf("Resolve(%q) should error", id)
		}
	}
}
