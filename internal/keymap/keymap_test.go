package keymap

import "testing"

func sample() *Keymap {
	return New(
		Layer{Scope: "global", Bindings: []Binding{
			{Action: "quit", Keys: []string{"q", "ctrl+c"}, Help: "quit"},
			{Action: "help", Keys: []string{"?"}, Help: "help"},
		}},
		Layer{Scope: "item", Bindings: []Binding{
			{Action: "complete", Keys: []string{"x", "space"}, Help: "complete"},
			{Action: "quit", Keys: []string{"Q"}, Help: "quit (item)"},
		}},
	)
}

func TestMatchPrefersSpecificScope(t *testing.T) {
	k := sample()
	active := []string{"item", "global"}

	if a, ok := k.Match("x", active); !ok || a != "complete" {
		t.Errorf("x -> %q,%v want complete", a, ok)
	}
	if a, ok := k.Match("q", active); !ok || a != "quit" {
		t.Errorf("q -> %q,%v want quit (inherited from global)", a, ok)
	}
	if a, ok := k.Match("Q", active); !ok || a != "quit" {
		t.Errorf("Q -> %q,%v want quit (item override)", a, ok)
	}
	if _, ok := k.Match("z", active); ok {
		t.Error("z should not match")
	}
	if _, ok := k.Match("x", []string{"global"}); ok {
		t.Error("x should not match when item scope is inactive")
	}
}

func TestHelpDedupMostSpecificWins(t *testing.T) {
	k := sample()
	help := k.Help([]string{"item", "global"})
	if len(help) != 3 {
		t.Fatalf("help len = %d: %+v", len(help), help)
	}
	if help[0].Action != "complete" || help[1].Action != "quit" || help[1].Help != "quit (item)" || help[2].Action != "help" {
		t.Errorf("help = %+v", help)
	}
}

func TestMergeAppliesAndReportsUnknown(t *testing.T) {
	k := sample()
	merged, unknown := k.Merge([]Override{
		{Scope: "global", Action: "quit", Keys: []string{"Z"}},
		{Scope: "item", Action: "nope", Keys: []string{"n"}},
	})
	if a, ok := merged.Match("Z", []string{"global"}); !ok || a != "quit" {
		t.Errorf("Z -> %q,%v want quit", a, ok)
	}
	if _, ok := k.Match("Z", []string{"global"}); ok {
		t.Error("original keymap must be unchanged by Merge")
	}
	if len(unknown) != 1 || unknown[0].Action != "nope" {
		t.Errorf("unknown = %+v", unknown)
	}
}

func TestOverridesFromMap(t *testing.T) {
	ovs := OverridesFromMap(map[string]map[string][]string{
		"global": {"quit": {"q", "ctrl+c"}},
	})
	if len(ovs) != 1 || ovs[0].Scope != "global" || ovs[0].Action != "quit" || len(ovs[0].Keys) != 2 {
		t.Errorf("OverridesFromMap = %+v", ovs)
	}
}
