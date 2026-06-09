package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	gotoml "github.com/pelletier/go-toml/v2"
)

// Load reads the config at path, falling back to built-in defaults for any
// unspecified field. A missing file yields the defaults.
func Load(path string) (Config, error) {
	cfg := Default()
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return cfg.expanded(), nil
	}
	if err != nil {
		return Config{}, err
	}
	if err := gotoml.Unmarshal(b, &cfg); err != nil {
		return Config{}, err
	}
	cfg = cfg.expanded()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// EnsureConfig loads the config, writing the defaults to path first if the file
// does not exist yet. It reports whether the file was just created.
func EnsureConfig(path string) (Config, bool, error) {
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err := writeDefault(path); err != nil {
			return Config{}, false, err
		}
		cfg, err := Load(path)
		return cfg, true, err
	}
	cfg, err := Load(path)
	return cfg, false, err
}

// writeDefault serializes the built-in defaults to path.
func writeDefault(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := gotoml.Marshal(Default())
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// expanded resolves ~ and environment variables in path fields.
func (c Config) expanded() Config {
	c.Store = expandPath(c.Store)
	c.Mirror = expandPath(c.Mirror)
	return c
}

// expandPath expands a leading ~ and any environment variables.
func expandPath(p string) string {
	if p == "" {
		return ""
	}
	if p == "~" || strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return os.ExpandEnv(p)
}
