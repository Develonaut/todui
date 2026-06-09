# Architecture — Layered, Interface-Bounded, DAG

> Adapted for Go from the bnto project's `architecture.md`. Each layer depends only on layers below it (toward abstractions), never up or sideways into concretes. Never skip layers.

## Dependency direction (inward only)

```
cmd/todui              composition root — builds concretes, wires interfaces, dispatches
   │
internal/app           use-cases / the single write path; depends on interfaces only
   │
internal/ports         the interface seams:  Repository · Codec · Renderer · Clock
   │            ▲
   │   implemented by adapters (chosen only at the root):
   │   store/file · codec/toml · render/markdown · clock/system
   │
internal/todo          PURE DOMAIN — imports nothing of ours, no I/O. Everyone depends on it; it depends on no one.
```

Drivers (`internal/cli`, `internal/tui`) sit at the top alongside the root's wiring: they translate user intent into `app` use-case calls. They never touch adapters or the filesystem directly.

## Package responsibilities

| Package | Responsibility | Depends on |
|---|---|---|
| `internal/todo` | Domain types, positional IDs, pure state transitions (add/edit/move/reorder/claim/trim), normalization | nothing (stdlib only) |
| `internal/ports` | Interface definitions (`Repository`, `Codec`, `Renderer`, `Clock`) | `todo` |
| `internal/config` | Generic, user-facing configuration (sections, store path, mirror, theme) + defaults + first-run bootstrap | `todo` |
| `internal/app` | Use-cases (one method per operation), each via `Repository.Update` | `ports`, `todo` |
| `internal/codec/toml` | `Codec`: `todo.List` ⇄ TOML bytes | `ports`, `todo` |
| `internal/render/markdown` | `Renderer`: `todo.List` → Markdown mirror bytes | `ports`, `todo`, `config` |
| `internal/store/file` | `Repository`: atomic, locked load/update; writes store + mirror | `ports`, `todo` |
| `internal/importer/markdown` | One-shot legacy-Markdown → `todo.List` | `todo`, `config` |
| `internal/clock` | `Clock`: system time (tests use a fake) | `ports` |
| `internal/cli` | CLI driver: `Command` interface + registry, `--json` output | `app`, `config`, `todo` |
| `internal/tui` | Bubble Tea v2 driver: lists, Huh forms, fsnotify reload | `app`, `config`, `todo` |
| `cmd/todui` | Composition root: the **only** place that constructs concrete adapters | everything |

## The interface seams (`internal/ports`)

```go
type Codec interface {
    Encode(todo.List) ([]byte, error)
    Decode([]byte) (todo.List, error)
}
type Renderer interface {
    Render(todo.List) ([]byte, error)
}
type Repository interface {
    Load() (todo.List, error)
    Update(func(*todo.List) error) error // lock → load → mutate → normalize/trim → save (+mirror) → unlock
}
type Clock interface{ Now() time.Time }
```

The CLI `Command` interface lives in `internal/cli` (it's a driver concern), self-describing and dispatched via a registry (a table lookup, **not** a giant switch).

## Hard rules
- **`internal/todo` stays pure.** A PR that makes it import `os`, `time` (in logic), or any of our other packages is wrong.
- **All mutations go through `Repository.Update`** with an intent closure keyed by item ID. Never save a whole cached list — that's how lost updates happen between the TUI and a concurrent CLI/skill invocation.
- **Concretes only at the root.** If a package outside `cmd/todui` constructs a `*toml.Codec` or opens a file, reconsider.
- **DAG only.** `go list` / the linter must show no import cycles.
