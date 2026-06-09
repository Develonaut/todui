# Testing — TDD Red-First

> Adapted for Go from bnto's testing discipline. Tests define the contract. Write the failing test first, watch it fail for the right reason, then implement.

## Workflow
1. **Red** — write a failing `_test.go` that expresses the desired behavior.
2. **Green** — implement the minimum to pass.
3. **Refactor** — clean up under green tests.

## Rules
- **Colocate tests.** `foo_test.go` lives next to `foo.go` in the same package (use an external `_test` package only when testing the exported surface specifically).
- **Table-driven tests** for domain logic (`internal/todo`): IDs, resolve, mutate, move/reorder, trim, normalize — cover happy path, boundaries, and error/ambiguous cases.
- **Golden tests** for the Markdown renderer and the importer: compare output against `testdata/golden/*`. Regenerate intentionally, never blindly.
- **Round-trip gate (acceptance):** `import(legacy TODO.md)` → `render` must equal the original (ignoring the `Last updated` line). This is the strongest single test that the migration is faithful.
- **Race detector** for anything touching the store: `go test -race`. Include a concurrent-`Update` test asserting no lost updates / no corruption.
- **Deterministic time:** inject a fake `Clock` in tests; never call the system clock from a test assertion path.
- **No skipped tests** on the main branch. `t.Skip` only for genuine environment gaps, with a reason.

## What to test where
- `internal/todo` — pure unit tests, no I/O, fast.
- `internal/store/file` — temp-dir (`t.TempDir()`) integration tests, atomic write + lock + `-race`.
- `internal/render/markdown`, `internal/importer/markdown` — golden + round-trip.
- `internal/config` — defaults, load precedence, first-run bootstrap.
- CLI/TUI — exercise via the `app` layer and command handlers where practical; keep terminal-driven paths thin.
