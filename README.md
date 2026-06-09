# todui

A tiny, config-driven task manager for your terminal — a CLI-first core with a
[Charm](https://charm.sh) Bubble Tea TUI layered on top. Tasks live in a simple
TOML file; todui can also render a read-only Markdown mirror for at-a-glance
reading.

> **Status:** early development. See [the milestones](#status) below.

## Why

Most terminal to-do tools are either a single opaque binary or a pile of shell
aliases over a Markdown file. todui keeps your tasks in a plain, hand-eyeball-able
TOML file, gives you a clean keyboard-driven TUI, and exposes a fully scriptable
CLI underneath — so the same data is equally good for a human and for automation.

## Install

```sh
go install github.com/develonaut/todui/cmd/todui@latest
```

This installs the `todui` binary into `$(go env GOPATH)/bin`.

## Quick start

```sh
todui              # launch the TUI (self-initializes config on first run)
todui add "Buy milk" --section now
todui list
todui done NB
```

On first run, todui writes a default config to `~/.config/todui/config.toml`
and creates an empty task store — no setup required.

## Configuration

Everything user-specific is configuration, not code. `~/.config/todui/config.toml`:

```toml
store  = "~/.local/share/todui/todo.toml"   # where tasks live
# mirror = "~/notes/TODO.md"                 # optional read-only Markdown view

[[section]]  key = "now"    title = "Now"    letter = "N"
[[section]]  key = "next"   title = "Next"   letter = "X"
[[section]]  key = "later"  title = "Later"  letter = "L"
[[section]]  key = "done"   title = "Done"   done = true

[done]
max_items = 10
max_age_days = 7
```

Sections, the store path, the optional Markdown mirror, and the theme are all
yours to define. Task IDs (like `NA`, `XB`) are positional and computed for
display — you never manage them by hand.

## CLI

```
todui [tui]                       launch the TUI (default)
todui list [--section S] [--json]
todui add <task...> [--section S] [--tag T]... [--ado REF] [--context C]
todui done   <id|query> [--date YYYY-MM-DD]
todui start  <id|query>           move to in-progress (claim)
todui mv     <id> <section>
todui reorder <id> up|down
todui edit   <id> [--task ...] [--context ...] [--tag ...]...
todui rm     <id|query> --yes
todui init   [--force] [--import FILE]
todui import [--from FILE] [--force]
todui render [--out FILE]
```

Every command supports `--json` for scripting.

## Keybindings (TUI)

`j`/`k` or arrows move · `Tab`/`<`/`>` switch section · `space`/`x` complete ·
`a` add · `e` edit · `d` delete · `s` start/claim · `J`/`K` reorder ·
`H`/`L` move section · `r` reload · `?` help · `q` quit.

## Contributing

This project follows a small set of engineering standards (single-responsibility
packages, interface-bounded layers, TDD). Please read
[`.claude/rules/`](./.claude/rules) before contributing.

```sh
make build   # compile
make test    # go test -race ./...
make lint    # golangci-lint + go vet
```

## Status

- [x] **M0** — skeleton & standards
- [ ] **M1** — domain core
- [ ] **M2** — ports + adapters (TOML, store, Markdown renderer)
- [ ] **M3** — app service + CLI
- [ ] **M4** — Markdown importer
- [ ] **M5** — TUI
- [ ] **M6** — live reload
- [ ] **M7** — release polish

## License

[MIT](./LICENSE)
