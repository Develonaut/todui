// Package keymap is a small, application-agnostic contextual keybinding system.
//
// Bindings are organized into named scopes (e.g. "global", "item", "form"). At
// any moment the application supplies an ordered "active" stack of scopes, most
// specific first. Both key dispatch (Match) and help generation (Help) read
// from this one source, so the help a user sees can never drift out of sync
// with what the keys actually do.
package keymap

import "slices"

// Binding maps one or more keys to an action, with help text for display.
type Binding struct {
	Action string
	Keys   []string
	Help   string
}

// Layer is the set of bindings for one scope.
type Layer struct {
	Scope    string
	Bindings []Binding
}

// Keymap holds bindings grouped by scope.
type Keymap struct {
	order  []string
	scopes map[string][]Binding
}

// New builds a Keymap from layers. Layers sharing a scope name are concatenated.
func New(layers ...Layer) *Keymap {
	k := &Keymap{scopes: make(map[string][]Binding)}
	for _, l := range layers {
		if _, ok := k.scopes[l.Scope]; !ok {
			k.order = append(k.order, l.Scope)
		}
		k.scopes[l.Scope] = append(k.scopes[l.Scope], l.Bindings...)
	}
	return k
}

// Match returns the action bound to key in the first active scope that defines
// it. Active scopes are searched in order (most specific first), so a context
// can override a global key while inheriting everything it does not redefine.
func (k *Keymap) Match(key string, active []string) (string, bool) {
	for _, scope := range active {
		for _, b := range k.scopes[scope] {
			if slices.Contains(b.Keys, key) {
				return b.Action, true
			}
		}
	}
	return "", false
}

// Help returns the bindings visible in the active scopes, most specific first,
// with each action appearing once (the most specific binding wins). This is the
// single source the help bar should render.
func (k *Keymap) Help(active []string) []Binding {
	seen := make(map[string]bool)
	var out []Binding
	for _, scope := range active {
		for _, b := range k.scopes[scope] {
			if seen[b.Action] {
				continue
			}
			seen[b.Action] = true
			out = append(out, b)
		}
	}
	return out
}
