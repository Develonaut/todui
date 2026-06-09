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
	actFoldToggle  = "fold.toggle"
	actExpand      = "fold.expand"
	actCollapse    = "fold.collapse"
	actGoalUp      = "goal.up"
	actGoalDown    = "goal.down"
)

// Scope names for the active stack.
const (
	scopeGlobal  = "global"
	scopeList    = "list"
	scopeItem    = "item"
	scopeHeader  = "header"
	scopeEmpty   = "empty"
	scopeConfirm = "confirm"
)

// defaultKeymap is the built-in binding set. Arrow keys are the advertised
// controls; vim keys (j/k/h/l) and the up/down halves of paired actions are
// kept as Hidden aliases so they still work but don't clutter the help bar.
// Paired actions use HelpKey to show a single combined hint (e.g. "↑↓").
func defaultKeymap() *keymap.Keymap {
	return keymap.New(
		keymap.Layer{Scope: scopeGlobal, Bindings: []keymap.Binding{
			{Action: actGoalUp, Keys: []string{"+", "="}, Hidden: true},
			{Action: actGoalDown, Keys: []string{"-", "_"}, Hidden: true},
			{Action: actHelp, Keys: []string{"?"}, Help: "help"},
			{Action: actReload, Keys: []string{"r"}, Help: "reload", Hidden: true},
			{Action: actQuit, Keys: []string{"q", "ctrl+c"}, Help: "quit"},
		}},
		keymap.Layer{Scope: scopeList, Bindings: []keymap.Binding{
			{Action: actUp, Keys: []string{"up", "k"}, Hidden: true},
			{Action: actDown, Keys: []string{"down", "j"}, Help: "navigate", HelpKey: "↑↓"},
			{Action: actExpand, Keys: []string{"right", "l"}, Help: "fold", HelpKey: "←→"},
			{Action: actCollapse, Keys: []string{"left", "h"}, Hidden: true},
			{Action: actSectionNext, Keys: []string{"tab"}, Hidden: true},
			{Action: actSectionPrev, Keys: []string{"shift+tab"}, Hidden: true},
			{Action: actAdd, Keys: []string{"a"}, Help: "add"},
		}},
		keymap.Layer{Scope: scopeItem, Bindings: []keymap.Binding{
			{Action: actComplete, Keys: []string{"d"}, Help: "done"},
			{Action: actEdit, Keys: []string{"e"}, Help: "edit"},
			{Action: actStart, Keys: []string{"s"}, Help: "start"},
			{Action: actDelete, Keys: []string{"x", "delete", "backspace"}, Help: "delete"},
			{Action: actReorderUp, Keys: []string{"shift+up", "K"}, Hidden: true},
			{Action: actReorderDown, Keys: []string{"shift+down", "J"}, Help: "reorder", HelpKey: "⇧↑↓"},
			{Action: actMovePrev, Keys: []string{"shift+left", "H"}, Hidden: true},
			{Action: actMoveNext, Keys: []string{"shift+right", "L"}, Help: "move", HelpKey: "⇧←→"},
		}},
		keymap.Layer{Scope: scopeHeader, Bindings: []keymap.Binding{
			{Action: actFoldToggle, Keys: []string{"enter"}, Hidden: true},
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
		if r, ok := m.selectedRow(); ok && r.header {
			return []string{scopeHeader, scopeList, scopeGlobal}
		}
		return []string{scopeItem, scopeList, scopeGlobal}
	default:
		return []string{scopeGlobal}
	}
}
