package todo

import "strings"

// Section describes one bucket of the list. Its values (key, title, letter,
// done flag) come from configuration; the domain defines only the shape and the
// identity/ordering rules that operate on it.
type Section struct {
	Key    string // stable identifier stored on items, e.g. "now"
	Title  string // display heading, e.g. "Now"
	Letter string // single uppercase letter used to build IDs; empty for done
	Done   bool   // the terminal section (completed items; no IDs)
}

// Schema is the ordered set of sections a list is organized into. Display order
// is slice order; an item's section is matched to the schema by Key.
type Schema struct {
	Sections []Section
}

// index returns the position of key in display order, or -1 if absent.
func (s Schema) index(key string) int {
	for i := range s.Sections {
		if s.Sections[i].Key == key {
			return i
		}
	}
	return -1
}

// byKey looks up a section by its key.
func (s Schema) byKey(key string) (Section, bool) {
	if i := s.index(key); i >= 0 {
		return s.Sections[i], true
	}
	return Section{}, false
}

// DoneKey returns the key of the terminal (done) section, or "" if none.
func (s Schema) DoneKey() string {
	for i := range s.Sections {
		if s.Sections[i].Done {
			return s.Sections[i].Key
		}
	}
	return ""
}

// rank orders sections for normalization; unknown sections sort to the end.
func (s Schema) rank(key string) int {
	if i := s.index(key); i >= 0 {
		return i
	}
	return len(s.Sections)
}

// splitID separates an uppercased ID into its section and sequence label,
// matching the longest section letter that prefixes it.
func (s Schema) splitID(id string) (Section, string, bool) {
	var best Section
	var rest string
	found := false
	for i := range s.Sections {
		sec := s.Sections[i]
		if sec.Letter == "" || !strings.HasPrefix(id, sec.Letter) {
			continue
		}
		if !found || len(sec.Letter) > len(best.Letter) {
			best, rest, found = sec, id[len(sec.Letter):], true
		}
	}
	return best, rest, found
}
