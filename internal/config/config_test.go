package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultIsValid(t *testing.T) {
	if err := Default().Validate(); err != nil {
		t.Fatal(err)
	}
	if got := Default().DefaultSection(); got != "now" {
		t.Errorf("DefaultSection() = %q, want now", got)
	}
}

func TestValidateRejectsBadConfigs(t *testing.T) {
	cases := map[string]Config{
		"no sections":    {Sections: nil},
		"no done":        {Sections: []SectionDef{{Key: "a", Letter: "A"}}},
		"missing letter": {Sections: []SectionDef{{Key: "a"}, {Key: "d", Done: true}}},
		"bad start":      {Sections: []SectionDef{{Key: "a", Letter: "A"}, {Key: "d", Done: true}}, StartSection: "x"},
	}
	for name, c := range cases {
		if err := c.Validate(); err == nil {
			t.Errorf("%s: expected invalid", name)
		}
	}
}

func TestLoadMissingReturnsDefaults(t *testing.T) {
	cfg, err := Load(filepath.Join(t.TempDir(), "nope.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Sections) != 4 {
		t.Errorf("want 4 default sections, got %d", len(cfg.Sections))
	}
}

func TestEnsureConfigCreatesThenLoads(t *testing.T) {
	p := filepath.Join(t.TempDir(), "config.toml")
	cfg, created, err := EnsureConfig(p)
	if err != nil || !created {
		t.Fatalf("created=%v err=%v", created, err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Errorf("config file not written: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("written config invalid: %v", err)
	}
	if _, created2, _ := EnsureConfig(p); created2 {
		t.Error("second EnsureConfig should not recreate")
	}
}

func TestExpandPathTilde(t *testing.T) {
	home, _ := os.UserHomeDir()
	if got := expandPath("~/x"); got != filepath.Join(home, "x") {
		t.Errorf("expandPath(~/x) = %q", got)
	}
}
