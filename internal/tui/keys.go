package tui

import "github.com/develonaut/todui/internal/keymap"

// Action identifiers. The keymap binds keys to these; dispatch maps these to
// behavior. Keeping them as named constants avoids stringly-typed drift.
const (
	actUp          = "cursor.up"
	actDown        = "cursor.down"
	actSectionPrev = "section.prev"
	actSectionNext = "section.next"
	actComplete    = "complete"
	actAdd         = "add"
	actEdit        = "edit"
	actDelete      = "delete"
	actStart       = "start"
	actReorderUp   = "reorder.up"
	actReorderDown = "reorder.down"
	actMovePrev    = "move.prev"
	actMoveNext    = "move.next"
	actReload      = "reload"
	actHelp        = "help"
	actQuit        = "quit"
	actConfirmYes  = "confirm.yes"
	actConfirmNo   = "confirm.no"
)

// Scope names for the active stack.
const (
	scopeGlobal  = "global"
	scopeList    = "list"
	scopeItem    = "item"
	scopeEmpty   = "empty"
	scopeConfirm = "confirm"
)

// defaultKeymap is the built-in binding set. Universal actions live in global;
// list-navigation in list; per-item actions in item; the rest in their scope.
// Bindings are generous (e.g. both "space" and " ", both "J" and "shift+j") so
// they work regardless of how the terminal reports the key.
func defaultKeymap() *keymap.Keymap {
	return keymap.New(
		keymap.Layer{Scope: scopeGlobal, Bindings: []keymap.Binding{
			{Action: actHelp, Keys: []string{"?"}, Help: "help"},
			{Action: actReload, Keys: []string{"r"}, Help: "reload"},
			{Action: actQuit, Keys: []string{"q", "ctrl+c"}, Help: "quit"},
		}},
		keymap.Layer{Scope: scopeList, Bindings: []keymap.Binding{
			{Action: actUp, Keys: []string{"k", "up"}, Help: "up"},
			{Action: actDown, Keys: []string{"j", "down"}, Help: "down"},
			{Action: actSectionPrev, Keys: []string{"shift+tab", "h", "left", "<"}, Help: "prev section"},
			{Action: actSectionNext, Keys: []string{"tab", "l", "right", ">"}, Help: "next section"},
			{Action: actAdd, Keys: []string{"a"}, Help: "add"},
		}},
		keymap.Layer{Scope: scopeItem, Bindings: []keymap.Binding{
			{Action: actComplete, Keys: []string{"space", " ", "x"}, Help: "done"},
			{Action: actEdit, Keys: []string{"e"}, Help: "edit"},
			{Action: actStart, Keys: []string{"s"}, Help: "start"},
			{Action: actDelete, Keys: []string{"d"}, Help: "delete"},
			{Action: actReorderUp, Keys: []string{"K", "shift+up"}, Help: "move up"},
			{Action: actReorderDown, Keys: []string{"J", "shift+down"}, Help: "move down"},
			{Action: actMovePrev, Keys: []string{"H", "shift+left"}, Help: "← section"},
			{Action: actMoveNext, Keys: []string{"L", "shift+right"}, Help: "section →"},
		}},
		keymap.Layer{Scope: scopeEmpty, Bindings: nil},
		keymap.Layer{Scope: scopeConfirm, Bindings: []keymap.Binding{
			{Action: actConfirmYes, Keys: []string{"y"}, Help: "yes"},
			{Action: actConfirmNo, Keys: []string{"n", "esc"}, Help: "no"},
		}},
	)
}

// activeScopes returns the scope stack for the current mode, most specific first.
func (m *Model) activeScopes() []string {
	switch m.mode {
	case modeConfirm:
		return []string{scopeConfirm, scopeGlobal}
	case modeList:
		if len(m.rows) == 0 {
			return []string{scopeEmpty, scopeList, scopeGlobal}
		}
		return []string{scopeItem, scopeList, scopeGlobal}
	default:
		return []string{scopeGlobal}
	}
}
