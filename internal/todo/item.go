// Package todo is the pure domain core: task types and the rules for
// organizing, identifying, ordering, and trimming them. It has no I/O and
// depends on nothing outside the standard library.
package todo

// Item is a single task. Field tags describe its on-disk (TOML) shape, but the
// domain treats Item as a plain value — persistence lives in adapter packages.
//
// Positional identifiers (e.g. "NA") are never stored; they are computed from
// the item's section and order at display time. Tags hold bare slugs
// ("mtb/ui"); brackets and backticks are presentation. ADO holds only an
// optional leading reference token ("#378701"); any other inline references
// stay within Task verbatim.
type Item struct {
	Task     string   `toml:"task"`
	Context  string   `toml:"context,omitempty"`
	Tags     []string `toml:"tags,omitempty"`
	ADO      string   `toml:"ado,omitempty"`
	Claimed  bool     `toml:"claimed,omitempty"`
	Section  string   `toml:"section"`
	Order    int      `toml:"order"`
	DoneDate string   `toml:"done_date,omitempty"`
}
