package todo

import (
	"fmt"
	"strings"
)

// seqLabel converts a 0-based index to a bijective base-26 label:
// 0->"A", 25->"Z", 26->"AA", 27->"AB".
func seqLabel(n int) string {
	n++ // shift to 1-indexed for bijective base-26
	var b []byte
	for n > 0 {
		n--
		b = append([]byte{byte('A' + n%26)}, b...)
		n /= 26
	}
	return string(b)
}

// seqIndex is the inverse of seqLabel. It reports ok=false for any input that
// is not a non-empty run of A-Z.
func seqIndex(s string) (int, bool) {
	if s == "" {
		return 0, false
	}
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			return 0, false
		}
		n = n*26 + int(c-'A'+1)
	}
	return n - 1, true
}

// ID returns the positional identifier for an item in section sec at the given
// normalized order (e.g. "NA"), or "" for the done section or a section without
// a letter.
func (s Schema) ID(sec Section, order int) string {
	if sec.Done || sec.Letter == "" {
		return ""
	}
	return sec.Letter + seqLabel(order)
}

// ComputeID returns the positional identifier for the item at idx (e.g. "NA"),
// or "" for items in the done section. It assumes the list is normalized so
// that an item's Order is its position within its section.
func (l *List) ComputeID(s Schema, idx int) string {
	it := l.Items[idx]
	sec, ok := s.byKey(it.Section)
	if !ok {
		return ""
	}
	return s.ID(sec, it.Order)
}

// Resolve maps a positional identifier back to its item index. It is
// case-insensitive and tolerates surrounding whitespace. It assumes the list is
// normalized.
func (l *List) Resolve(s Schema, id string) (int, error) {
	id = strings.ToUpper(strings.TrimSpace(id))
	sec, rest, ok := s.splitID(id)
	if !ok {
		return -1, fmt.Errorf("todo: unknown id %q", id)
	}
	seq, ok := seqIndex(rest)
	if !ok {
		return -1, fmt.Errorf("todo: malformed id %q", id)
	}
	for i := range l.Items {
		if l.Items[i].Section == sec.Key && l.Items[i].Order == seq {
			return i, nil
		}
	}
	return -1, fmt.Errorf("todo: no item with id %q", id)
}
