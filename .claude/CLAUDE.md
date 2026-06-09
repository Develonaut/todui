# todui — Agent & Contributor Instructions

**todui** is a tiny, open-source, config-driven task manager for the terminal: a CLI-first core with a [Charm](https://charm.sh) Bubble Tea TUI layered on top. It stores tasks in a TOML file and can render a read-only Markdown mirror.

## Read before any work

This project follows a small set of **engineering standards** adapted from the "Bento Box" rules of the [bnto](https://github.com/) project. **Review `.claude/rules/` before starting, and confirm you have, before writing code.**

- [`rules/code-standards.md`](./rules/code-standards.md) — the Bento Box Principle (single responsibility, no grab-bag files, small functions, YAGNI, comments).
- [`rules/architecture.md`](./rules/architecture.md) — layered + dependency-inverted design, the interface seams, package responsibilities.
- [`rules/testing.md`](./rules/testing.md) — TDD red-first, colocated tests, golden + round-trip gates.
- [`rules/pre-commit.md`](./rules/pre-commit.md) — the mandatory quality gate.
- [`rules/gotchas.md`](./rules/gotchas.md) — known pitfalls (Charm v2, fsnotify, file locking).

## The one-paragraph architecture

`internal/todo` is the **pure domain** — types, positional-ID computation, and state-transition operations, with **zero I/O and zero dependencies on our other packages**. Everything else depends on it; it depends on nothing. I/O sits behind interfaces in `internal/ports` (`Repository`, `Codec`, `Renderer`, `Clock`); concrete adapters (`codec/toml`, `render/markdown`, `store/file`, `clock`) are chosen **only at the composition root** (`cmd/todui/main.go`). The CLI (`internal/cli`) and TUI (`internal/tui`) are two thin drivers over a single use-case layer (`internal/app`), so both go through one write path and can never corrupt the store.

## Common commands

```sh
make build     # go build ./...
make test      # go test -race ./...
make lint      # golangci-lint run + go vet
make install   # go install ./cmd/todui  (+ td symlink)
make run       # go run ./cmd/todui tui
```

## Core values

- Single responsibility per package/file. **No `utils`/`helpers`/`common` grab bags.**
- Domain logic is pure; I/O, time, and the terminal are injected behind interfaces.
- Tests come first (red), and are colocated with the code they test.
- Nothing user-specific is hardcoded — sections, paths, mirror, and theme are configuration.
- Keep it tiny. YAGNI. Don't add an abstraction until a real second case earns it.
