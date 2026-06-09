package config

import (
	"os"
	"path/filepath"
)

// dataDir is where the task store lives by default (XDG data home).
func dataDir() string {
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, "todui")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "todui")
}

// configDir is where the config file lives by default (XDG config home).
func configDir() string {
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, "todui")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "todui")
}

// DefaultPath is the default config file location.
func DefaultPath() string {
	return filepath.Join(configDir(), "config.toml")
}

// Default returns the built-in configuration used when no config file exists.
// The default flow is Now / Next / Later / Done.
func Default() Config {
	return Config{
		Store: filepath.Join(dataDir(), "todo.toml"),
		Sections: []SectionDef{
			{Key: "now", Title: "Now", Letter: "N"},
			{Key: "next", Title: "Next", Letter: "X"},
			{Key: "later", Title: "Later", Letter: "L"},
			{Key: "done", Title: "Done", Letter: "D", Done: true},
		},
		Done: DoneConfig{MaxItems: 10, MaxAgeDays: 7},
	}
}
