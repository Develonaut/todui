package todo

import "sort"

// List is the whole task collection plus the header lines and timestamp used
// when rendering a Markdown mirror. Items carry their own section and order;
// the list does not store positional IDs.
type List struct {
	Header      []string `toml:"header,omitempty"`
	LastUpdated string   `toml:"last_updated"`
	Items       []Item   `toml:"item"`
}

// Section returns the items in a section, in display (order) sequence. The
// returned slice is a copy of the pointers' values; mutate via the List.
func (l *List) Section(s Schema, key string) []Item {
	var out []Item
	for i := range l.Items {
		if l.Items[i].Section == key {
			out = append(out, l.Items[i])
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Order < out[j].Order })
	return out
}

// Normalize sorts items by section (schema order) then by their current order,
// and renumbers each section's order to a dense 0..n-1 sequence. After
// Normalize, an item's Order equals its position within its section, which is
// what ID computation and resolution rely on.
func (l *List) Normalize(s Schema) {
	sort.SliceStable(l.Items, func(i, j int) bool {
		ri, rj := s.rank(l.Items[i].Section), s.rank(l.Items[j].Section)
		if ri != rj {
			return ri < rj
		}
		return l.Items[i].Order < l.Items[j].Order
	})
	counts := make(map[string]int)
	for i := range l.Items {
		key := l.Items[i].Section
		l.Items[i].Order = counts[key]
		counts[key]++
	}
}
