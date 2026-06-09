# Pre-Commit Quality Gate

> Adapted for Go from bnto's `pre-commit.md`. Run this gate before every commit. All steps must pass — no exceptions, no skipped tests.

```sh
go build ./...            # compiles
go vet ./...              # vet clean
golangci-lint run         # lint clean
go test -race ./...       # all tests pass, race-free
gofmt -l .                # prints nothing (everything formatted)
```

`make lint && make test` runs the core of this.

## Checklist
- [ ] Code is `gofmt`-formatted (`gofmt -l .` is empty).
- [ ] `go vet` and `golangci-lint run` are clean (fix, don't `//nolint` away — justify any rare exception inline).
- [ ] `go test -race ./...` is green; new behavior has a test written **before** the implementation.
- [ ] No new `utils`/`helpers`/`common` grab-bag file or package.
- [ ] No new import cycle; `internal/todo` still imports nothing of ours and no I/O.
- [ ] Files stay focused (well under ~300 production lines); functions stay small.
- [ ] Exported identifiers have doc comments.
- [ ] Commit message is concise and explains *why*.
