package todo

import "fmt"

// Move relocates the item at idx to the target section, appended at the bottom.
func (l *List) Move(s Schema, idx int, toKey string) error {
	if err := l.inRange(idx); err != nil {
		return err
	}
	if _, ok := s.byKey(toKey); !ok {
		return fmt.Errorf("todo: unknown section %q", toKey)
	}
	l.Items[idx].Section = toKey
	l.Items[idx].Order = l.nextOrder(toKey, idx)
	return nil
}

// Reorder moves the item at idx within its section by one slot toward the top
// (delta<0) or bottom (delta>0). At a section boundary it is a no-op.
func (l *List) Reorder(idx, delta int) error {
	if err := l.inRange(idx); err != nil {
		return err
	}
	if delta == 0 {
		return nil
	}
	step := 1
	if delta < 0 {
		step = -1
	}
	target := l.Items[idx].Order + step
	for i := range l.Items {
		if l.Items[i].Section == l.Items[idx].Section && l.Items[i].Order == target {
			l.Items[i].Order, l.Items[idx].Order = l.Items[idx].Order, l.Items[i].Order
			return nil
		}
	}
	return nil
}

// Complete moves the item at idx into the done section and stamps the date.
func (l *List) Complete(s Schema, idx int, date string) error {
	if err := l.inRange(idx); err != nil {
		return err
	}
	doneKey := s.DoneKey()
	if doneKey == "" {
		return fmt.Errorf("todo: schema has no done section")
	}
	l.Items[idx].Section = doneKey
	l.Items[idx].DoneDate = date
	l.Items[idx].Order = 0
	return nil
}
