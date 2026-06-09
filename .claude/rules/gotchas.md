# Gotchas — Known Pitfalls

> Captured as we hit them, so we don't relearn them. Add to this file when a non-obvious trap costs you time.

## Charm v2 (Bubble Tea / Bubbles / Lip Gloss / Huh)
- **Modules live under `charm.land/...`**, not `github.com/charmbracelet/...`. Use `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2`, `charm.land/huh/v2`. Mixing the old github paths with the new ones confuses `go mod`.
- **`Model.View()` returns `tea.View`, not `string`.** Child components (Bubbles `list`, Huh, Lip Gloss) produce strings. Compose them into one string, then wrap **once**: `return tea.NewView(finalString)`.
- **Alt-screen is a `tea.View` field**, not a program option — there's no `WithAltScreen`. Set it on the returned view. Verify the exact field name against the installed `bubbletea/v2` version.
- **Huh `Form.Update` returns `(huh.Model, tea.Cmd)`**, not `(tea.Model, ...)`. Type-assert back to `*huh.Form` before re-storing it. Drive the form to completion via `form.State == huh.StateCompleted` (and handle `StateAborted`).
- Key messages are `tea.KeyPressMsg` in v2.

## fsnotify
- **Watch the directory, not the file.** Atomic writes rename a temp file over the target, which replaces the inode — a watch on the file path stops firing after the first rename. Watch the parent dir and filter events by base name.
- Our own writes also fire the watcher; a redundant reload of a tiny file is harmless, so don't over-engineer debouncing.

## File locking & atomic writes
- Lock a **sidecar** file (`todo.toml.lock`), never the data file you're about to `rename` over. Use `lockedfile.MutexAt` (from `github.com/rogpeppe/go-internal/lockedfile`) — Go-toolchain-grade, preferred over the unmaintained `gofrs/flock`.
- Write the temp file in the **same directory** as the target so `os.Rename` is atomic (same filesystem). `fsync` the temp file before rename.
- The Markdown mirror is written `0444` (read-only) to signal "don't hand-edit — todui owns this."
