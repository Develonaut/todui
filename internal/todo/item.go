// Package todo is the pure domain core: task types and the rules for
// organizing, identifying, ordering, and trimming them. It has no I/O and
// depends on nothing outside the standard library.
package todo

import "strings"

// Item is a single task. Field tags describe its on-disk (TOML) shape, but the
// domain treats Item as a plain value — persistence lives in adapter packages.
//
// Title is a short headline (a few words); Description is the fuller context.
// Positional identifiers (e.g. "NA") are never stored; they are computed from
// the item's section and order at display time. Tags hold bare slugs
// ("mtb/ui"); brackets and backticks are presentation. ADO holds only an
// optional leading reference token ("#378701"); other inline references stay
// within the text verbatim.
type Item struct {
	Title       string   `toml:"title"`
	Description string   `toml:"description,omitempty"`
	Tags        []string `toml:"tags,omitempty"`
	ADO         string   `toml:"ado,omitempty"`
	Section     string   `toml:"section"`
	Order       int      `toml:"order"`
	DoneDate    string   `toml:"done_date,omitempty"`
}

// TitleSep separates a title from its description in rendered/parsed text.
const TitleSep = " — "

// SplitTitle splits free text into a short title and a fuller description at the
// first TitleSep. Text without the separator is all title.
func SplitTitle(text string) (title, desc string) {
	if i := strings.Index(text, TitleSep); i >= 0 {
		return text[:i], text[i+len(TitleSep):]
	}
	return text, ""
}

// JoinTitle is the inverse of SplitTitle.
func JoinTitle(title, desc string) string {
	if desc == "" {
		return title
	}
	return title + TitleSep + desc
}
