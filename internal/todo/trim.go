package todo

import (
	"sort"
	"time"
)

// doneDateLayout is the date format stored on completed items.
const doneDateLayout = "2006-01-02"

// TrimDone keeps only the most recent completed items: at most maxItems (when
// maxItems>0) and, when maxAgeDays>0, only those completed within that many days
// of now. Kept done items are renumbered newest-first; open items are untouched.
func (l *List) TrimDone(s Schema, now time.Time, maxItems, maxAgeDays int) error {
	doneKey := s.DoneKey()
	if doneKey == "" {
		return nil
	}

	var done, open []Item
	for _, it := range l.Items {
		if it.Section == doneKey {
			done = append(done, it)
		} else {
			open = append(open, it)
		}
	}

	sort.SliceStable(done, func(i, j int) bool {
		return doneTime(done[i].DoneDate).After(doneTime(done[j].DoneDate))
	})

	cutoff := now.AddDate(0, 0, -maxAgeDays)
	var kept []Item
	for _, it := range done {
		if maxItems > 0 && len(kept) >= maxItems {
			break
		}
		if maxAgeDays > 0 {
			dt := doneTime(it.DoneDate)
			if dt.IsZero() || dt.Before(cutoff) {
				continue
			}
		}
		it.Order = len(kept)
		kept = append(kept, it)
	}

	l.Items = append(open, kept...)
	return nil
}

// doneTime parses the leading date of a completion annotation, returning the
// zero time if absent. DoneDate may carry a trailing note (e.g.
// "2026-06-08, /standup"), so only the date prefix is read.
func doneTime(s string) time.Time {
	if len(s) < len(doneDateLayout) {
		return time.Time{}
	}
	t, err := time.Parse(doneDateLayout, s[:len(doneDateLayout)])
	if err != nil {
		return time.Time{}
	}
	return t
}
