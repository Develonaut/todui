// Package ports defines the interface seams between the pure domain and the
// outside world. Concrete adapters implement these interfaces; the composition
// root selects and wires them. Nothing here performs I/O.
package ports

import (
	"time"

	"github.com/develonaut/todui/internal/todo"
)

// Codec converts a todo.List to and from its serialized byte form.
type Codec interface {
	Encode(todo.List) ([]byte, error)
	Decode([]byte) (todo.List, error)
}

// Renderer turns a todo.List into presentation bytes, such as a Markdown mirror.
type Renderer interface {
	Render(todo.List) ([]byte, error)
}

// Repository is the persistent task store. Update applies fn within a single
// locked load → mutate → save critical section so concurrent callers cannot
// lose each other's writes.
type Repository interface {
	Load() (todo.List, error)
	Update(func(*todo.List) error) error
}

// Clock reports the current time. It is injected so that time-dependent logic
// stays deterministic under test.
type Clock interface {
	Now() time.Time
}
