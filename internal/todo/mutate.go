package todo

import "fmt"

// Add appends an item to its section (item.Section must already be set) and
// returns its index. The item is ordered after existing items in that section.
func (l *List) Add(it Item) int {
	it.Order = l.nextOrder(it.Section, -1)
	l.Items = append(l.Items, it)
	return len(l.Items) - 1
}

// Edit applies mut to the item at idx.
func (l *List) Edit(idx int, mut func(*Item)) error {
	if err := l.inRange(idx); err != nil {
		return err
	}
	mut(&l.Items[idx])
	return nil
}

// Delete removes the item at idx.
func (l *List) Delete(idx int) error {
	if err := l.inRange(idx); err != nil {
		return err
	}
	l.Items = append(l.Items[:idx], l.Items[idx+1:]...)
	return nil
}

// nextOrder returns one past the highest Order currently used in a section,
// ignoring the item at skip (use -1 to ignore none).
func (l *List) nextOrder(section string, skip int) int {
	max := -1
	for i := range l.Items {
		if i == skip || l.Items[i].Section != section {
			continue
		}
		if l.Items[i].Order > max {
			max = l.Items[i].Order
		}
	}
	return max + 1
}

func (l *List) inRange(idx int) error {
	if idx < 0 || idx >= len(l.Items) {
		return fmt.Errorf("todo: item index %d out of range", idx)
	}
	return nil
}
