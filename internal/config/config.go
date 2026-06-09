// Package config defines todui's user-facing configuration — sections, the
// store and mirror paths, and trim limits — plus its defaults and first-run
// bootstrap. Everything user-specific lives here as data, not in code.
package config

import (
	"fmt"

	"github.com/develonaut/todui/internal/todo"
)

// Config is the on-disk configuration.
type Config struct {
	Store        string       `toml:"store"`
	Mirror       string       `toml:"mirror,omitempty"`
	StartSection string       `toml:"start_section,omitempty"`
	Header       []string     `toml:"header,omitempty"`
	Sections     []SectionDef `toml:"section"`
	Done         DoneConfig   `toml:"done"`
}

// SectionDef configures one section.
type SectionDef struct {
	Key    string `toml:"key"`
	Title  string `toml:"title"`
	Letter string `toml:"letter,omitempty"`
	Done   bool   `toml:"done,omitempty"`
}

// DoneConfig limits how many completed items are retained.
type DoneConfig struct {
	MaxItems   int `toml:"max_items"`
	MaxAgeDays int `toml:"max_age_days"`
}

// Schema builds the domain schema from the configured sections.
func (c Config) Schema() todo.Schema {
	secs := make([]todo.Section, len(c.Sections))
	for i, s := range c.Sections {
		secs[i] = todo.Section{Key: s.Key, Title: s.Title, Letter: s.Letter, Done: s.Done}
	}
	return todo.Schema{Sections: secs}
}

// DefaultSection is the key new items are added to when none is given: the
// first non-done section.
func (c Config) DefaultSection() string {
	for _, s := range c.Sections {
		if !s.Done {
			return s.Key
		}
	}
	return ""
}

// Validate reports whether the configuration is internally consistent.
func (c Config) Validate() error {
	if len(c.Sections) == 0 {
		return fmt.Errorf("config: no sections defined")
	}
	seen := make(map[string]bool)
	doneCount := 0
	for _, s := range c.Sections {
		switch {
		case s.Key == "":
			return fmt.Errorf("config: section with empty key")
		case seen[s.Key]:
			return fmt.Errorf("config: duplicate section key %q", s.Key)
		}
		seen[s.Key] = true
		if s.Done {
			doneCount++
		} else if s.Letter == "" {
			return fmt.Errorf("config: section %q needs a letter", s.Key)
		}
	}
	if doneCount != 1 {
		return fmt.Errorf("config: exactly one section must be marked done (got %d)", doneCount)
	}
	if c.StartSection != "" && !seen[c.StartSection] {
		return fmt.Errorf("config: start_section %q is not a defined section", c.StartSection)
	}
	return nil
}
