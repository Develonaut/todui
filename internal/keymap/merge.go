package keymap

// Override replaces the keys bound to an action within a scope.
type Override struct {
	Scope  string
	Action string
	Keys   []string
}

// Merge returns a new Keymap with the overrides applied, leaving the receiver
// unchanged. Overrides that do not match an existing (scope, action) binding are
// returned as "unknown" so the caller can warn rather than silently drop them.
func (k *Keymap) Merge(overrides []Override) (*Keymap, []Override) {
	out := &Keymap{
		order:  append([]string(nil), k.order...),
		scopes: make(map[string][]Binding, len(k.scopes)),
	}
	for scope, bs := range k.scopes {
		cp := make([]Binding, len(bs))
		copy(cp, bs)
		out.scopes[scope] = cp
	}

	var unknown []Override
	for _, o := range overrides {
		if !out.apply(o) {
			unknown = append(unknown, o)
		}
	}
	return out, unknown
}

// apply sets the keys for a single override, reporting whether it matched.
func (k *Keymap) apply(o Override) bool {
	for i := range k.scopes[o.Scope] {
		if k.scopes[o.Scope][i].Action == o.Action {
			k.scopes[o.Scope][i].Keys = append([]string(nil), o.Keys...)
			return true
		}
	}
	return false
}

// OverridesFromMap converts nested scope→action→keys configuration into a flat
// slice of overrides.
func OverridesFromMap(m map[string]map[string][]string) []Override {
	var out []Override
	for scope, actions := range m {
		for action, keys := range actions {
			out = append(out, Override{Scope: scope, Action: action, Keys: keys})
		}
	}
	return out
}
