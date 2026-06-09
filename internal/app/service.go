// Package app holds the use-case layer: one method per task operation, each
// applied through the repository's single locked write path. It depends only on
// the domain and the port interfaces.
package app

import (
	"github.com/develonaut/todui/internal/ports"
	"github.com/develonaut/todui/internal/todo"
)

const (
	stampLayout = "2006-01-02 15:04"
	dateLayout  = "2006-01-02"
)

// Settings carries the domain policy the service applies on every write: the
// section schema, the section "start" moves items into (optional), and the
// done-trim limits.
type Settings struct {
	Schema       todo.Schema
	StartSection string
	DoneMax      int
	DoneAge      int
}

// Service exposes task use-cases over a repository and a clock.
type Service struct {
	repo ports.Repository
	clk  ports.Clock
	set  Settings
}

// NewService builds a Service.
func NewService(repo ports.Repository, clk ports.Clock, set Settings) *Service {
	return &Service{repo: repo, clk: clk, set: set}
}

// Schema returns the section schema in use, for display and ID computation.
func (s *Service) Schema() todo.Schema { return s.set.Schema }

// List returns the current tasks, normalized so positional IDs are stable.
func (s *Service) List() (todo.List, error) {
	l, err := s.repo.Load()
	if err != nil {
		return todo.List{}, err
	}
	l.Normalize(s.set.Schema)
	return l, nil
}

// Touch rewrites the store (and mirror) without changing tasks, materializing
// the files on first run.
func (s *Service) Touch() error {
	return s.repo.Update(func(l *todo.List) error {
		s.finalize(l)
		return nil
	})
}

// Add inserts a new item and returns its computed display ID.
func (s *Service) Add(it todo.Item) (string, error) {
	var id string
	err := s.repo.Update(func(l *todo.List) error {
		l.Add(it)
		s.finalize(l)
		id = s.idOfLastInSection(l, it.Section)
		return nil
	})
	return id, err
}

// Complete marks the item done with today's date.
func (s *Service) Complete(id string) error {
	return s.mutate(id, func(l *todo.List, idx int) error {
		return l.Complete(s.set.Schema, idx, s.clk.Now().Format(dateLayout))
	})
}

// Start moves the item into the start section (when configured) and claims it.
func (s *Service) Start(id string) error {
	return s.mutate(id, func(l *todo.List, idx int) error {
		if s.set.StartSection != "" {
			if err := l.Move(s.set.Schema, idx, s.set.StartSection); err != nil {
				return err
			}
		}
		return l.SetClaimed(idx, true)
	})
}

// Move relocates the item to another section.
func (s *Service) Move(id, section string) error {
	return s.mutate(id, func(l *todo.List, idx int) error {
		return l.Move(s.set.Schema, idx, section)
	})
}

// Reorder shifts the item within its section toward the top (delta<0) or
// bottom (delta>0).
func (s *Service) Reorder(id string, delta int) error {
	return s.mutate(id, func(l *todo.List, idx int) error {
		return l.Reorder(idx, delta)
	})
}

// Edit applies field changes to an item.
func (s *Service) Edit(id string, mut func(*todo.Item)) error {
	return s.mutate(id, func(l *todo.List, idx int) error {
		return l.Edit(idx, mut)
	})
}

// Delete removes an item.
func (s *Service) Delete(id string) error {
	return s.mutate(id, func(l *todo.List, idx int) error {
		return l.Delete(idx)
	})
}

// mutate resolves id within a locked update, runs action, then finalizes.
func (s *Service) mutate(id string, action func(*todo.List, int) error) error {
	return s.repo.Update(func(l *todo.List) error {
		l.Normalize(s.set.Schema)
		idx, err := l.Resolve(s.set.Schema, id)
		if err != nil {
			return err
		}
		if err := action(l, idx); err != nil {
			return err
		}
		s.finalize(l)
		return nil
	})
}

// finalize normalizes, trims completed items, and stamps the update time. It is
// the single place domain policy is applied on a write.
func (s *Service) finalize(l *todo.List) {
	l.Normalize(s.set.Schema)
	_ = l.TrimDone(s.set.Schema, s.clk.Now(), s.set.DoneMax, s.set.DoneAge)
	l.LastUpdated = s.clk.Now().Format(stampLayout)
}

// idOfLastInSection returns the display ID of the final item in a section,
// which is where Add places a new item.
func (s *Service) idOfLastInSection(l *todo.List, section string) string {
	sec, ok := s.set.Schema.Lookup(section)
	if !ok {
		return ""
	}
	n := 0
	for i := range l.Items {
		if l.Items[i].Section == section {
			n++
		}
	}
	if n == 0 {
		return ""
	}
	return s.set.Schema.ID(sec, n-1)
}
