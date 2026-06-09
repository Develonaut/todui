# Code Standards — The Bento Box Principle (Go)

> Adapted for Go from the bnto project's `code-standards.md`. Like a Japanese bento box where every compartment has one purpose and holds carefully prepared items, this codebase is organized with the same intention: small, focused, composable pieces that fit together cleanly.

## The five principles

1. **Single Responsibility** — each package, file, and function does **one** thing well.
2. **No Grab Bags** — there is no `utils`, `helpers`, `common`, `misc`, or `shared` package or file. Anything reusable gets a *named* home that describes what it is (`codec`, `clock`, `id`), never a dumping ground.
3. **Clear Boundaries** — well-defined interfaces between layers; no circular dependencies; the package graph is a DAG.
4. **Composable** — small pieces that combine. Prefer composing focused functions over growing one big one.
5. **YAGNI** — don't add features, parameters, or abstractions "just in case." An interface earns its place only when a real second implementation, a real I/O seam, or a real testability need exists.

## Concrete rules

### Files & packages
- **Package = one responsibility**, grouped by domain, not by technical pattern. Good: `render/markdown`, `codec/toml`. Bad: `utils`, `models` (as a junk drawer).
- **One major type or concern per file.** Group multiple small types in a file only when they are a single composable unit (e.g. `Item` + `List`). When a file grows a second responsibility, split it.
- **Size discipline:** target files well under **~300 production lines** (excluding `_test.go`); if a file approaches that, it is almost certainly doing two things — split it into a folder with focused files.
- **Functions target < 20 lines.** Long functions are a smell; extract named helpers (in the *same* domain file, not a `utils.go`).

### Naming & exports
- Exported identifiers carry their package as context — avoid stutter (`todo.List`, not `todo.TodoList`).
- Keep the exported surface minimal. Unexport anything a consumer doesn't need.

### Abstraction & dependencies
- **Domain core (`internal/todo`) imports nothing of ours and no I/O libraries.** It is pure.
- Define interfaces at I/O and abstraction seams (`internal/ports`). Consumers depend on the **interface**; the concrete type is selected **only at the composition root**.
- **No import cycles.** If two packages need the same logic, extract a third, small, *named* package for it.
- Inject collaborators (time via `Clock`, storage via `Repository`). **No global mutable state**, no `time.Now()` buried in domain logic.

### Comments
- A 1–3 line **file header** stating the file's purpose.
- `//` doc comments on every exported identifier (Go convention; the linter enforces it where configured).
- Explain **why**, not **what**. Don't narrate self-evident code. Document non-obvious decisions, domain rules, and format specs.

## The smell test
If you're about to create `utils.go`, a 400-line file, a function that scrolls off the screen, or an interface with one implementation and no test/second-impl justification — stop and reshape it. Small, named, single-purpose pieces.
