// Package file implements ports.Repository backed by a TOML file, with atomic
// writes, cross-process advisory locking, and an optional Markdown mirror.
package file

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/lockedfile"

	"github.com/develonaut/todui/internal/ports"
	"github.com/develonaut/todui/internal/todo"
)

// Options configures a Repository. Codec is required; Renderer and Mirror are
// optional and enable the read-only Markdown mirror when both are set.
type Options struct {
	Path     string
	Codec    ports.Codec
	Renderer ports.Renderer
	Mirror   string
}

// Repository is a file-backed todo store.
type Repository struct {
	opts Options
}

// New returns a file-backed repository.
func New(opts Options) *Repository {
	return &Repository{opts: opts}
}

var _ ports.Repository = (*Repository)(nil)

// Load reads and decodes the store. A missing file yields an empty list rather
// than an error, so first-run callers see an empty list.
func (r *Repository) Load() (todo.List, error) {
	b, err := os.ReadFile(r.opts.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return todo.List{}, nil
	}
	if err != nil {
		return todo.List{}, err
	}
	return r.opts.Codec.Decode(b)
}

// Update applies fn within an exclusive lock: load → fn → save (+mirror). The
// lock is held only for the duration of the call.
func (r *Repository) Update(fn func(*todo.List) error) error {
	unlock, err := r.lock()
	if err != nil {
		return err
	}
	defer unlock()

	list, err := r.Load()
	if err != nil {
		return err
	}
	if err := fn(&list); err != nil {
		return err
	}
	return r.save(list)
}

// lock takes the sidecar advisory lock, creating the store directory if needed.
func (r *Repository) lock() (func(), error) {
	if err := os.MkdirAll(filepath.Dir(r.opts.Path), 0o755); err != nil {
		return nil, err
	}
	return lockedfile.MutexAt(r.opts.Path + ".lock").Lock()
}

// save encodes the list to the store and, when configured, renders the mirror.
func (r *Repository) save(l todo.List) error {
	data, err := r.opts.Codec.Encode(l)
	if err != nil {
		return err
	}
	if err := writeAtomic(r.opts.Path, data, 0o644); err != nil {
		return err
	}
	if r.opts.Renderer == nil || r.opts.Mirror == "" {
		return nil
	}
	md, err := r.opts.Renderer.Render(l)
	if err != nil {
		return err
	}
	return writeAtomic(r.opts.Mirror, md, 0o444)
}
