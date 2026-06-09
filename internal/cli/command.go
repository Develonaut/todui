// Package cli is the command-line driver: a small registry of subcommands over
// the application service. It performs no storage wiring itself — concrete
// adapters are constructed by the composition root and handed in via Context.
package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/develonaut/todui/internal/app"
	"github.com/develonaut/todui/internal/config"
	"github.com/develonaut/todui/internal/ports"
)

// Context is what every command receives.
type Context struct {
	Svc      *app.Service
	Cfg      config.Config
	Renderer ports.Renderer
	Out      io.Writer
	Err      io.Writer
	JSON     bool
}

// Command is one subcommand.
type Command interface {
	Name() string
	Usage() string
	Run(cx *Context, args []string) error
}

// GlobalOpts are flags parsed before the subcommand.
type GlobalOpts struct {
	ConfigPath string
	StorePath  string
}

// ParseGlobal extracts leading global flags (--config, --file) and returns the
// remaining arguments (subcommand and its flags).
func ParseGlobal(args []string) (GlobalOpts, []string) {
	opts := GlobalOpts{ConfigPath: config.DefaultPath()}
	i := 0
	for i < len(args) {
		a := args[i]
		switch {
		case a == "--config" || a == "-config":
			if i+1 < len(args) {
				opts.ConfigPath = args[i+1]
				i += 2
				continue
			}
		case strings.HasPrefix(a, "--config="):
			opts.ConfigPath = strings.TrimPrefix(a, "--config=")
			i++
			continue
		case a == "--file" || a == "-file":
			if i+1 < len(args) {
				opts.StorePath = args[i+1]
				i += 2
				continue
			}
		case strings.HasPrefix(a, "--file="):
			opts.StorePath = strings.TrimPrefix(a, "--file=")
			i++
			continue
		}
		break
	}
	return opts, args[i:]
}

// Registry maps command names to commands.
type Registry map[string]Command

// NewRegistry builds the set of available commands.
func NewRegistry() Registry {
	cmds := []Command{
		listCmd{}, addCmd{}, doneCmd{}, startCmd{}, mvCmd{}, reorderCmd{},
		editCmd{}, rmCmd{}, renderCmd{}, resolveCmd{}, importCmd{}, initCmd{}, tuiCmd{},
	}
	r := make(Registry, len(cmds))
	for _, c := range cmds {
		r[c.Name()] = c
	}
	return r
}

// Dispatch runs the named command (defaulting to the TUI) and returns an exit
// code.
func (r Registry) Dispatch(cx *Context, args []string) int {
	if len(args) == 0 {
		args = []string{"tui"}
	}
	name, rest := args[0], args[1:]
	if name == "help" || name == "-h" || name == "--help" {
		r.printUsage(cx.Out)
		return 0
	}
	cmd, ok := r[name]
	if !ok {
		fmt.Fprintf(cx.Err, "todui: unknown command %q\n\n", name)
		r.printUsage(cx.Err)
		return 2
	}
	if err := cmd.Run(cx, rest); err != nil {
		if cx.JSON {
			_ = emitJSON(cx.Out, false, nil, nil, err)
		} else {
			fmt.Fprintln(cx.Err, "todui:", err)
		}
		return 1
	}
	return 0
}

// printUsage lists the available commands.
func (r Registry) printUsage(w io.Writer) {
	fmt.Fprintln(w, "todui — a tiny task manager\n\nUsage: todui [--config PATH] [--file STORE] <command> [args]\n\nCommands:")
	names := make([]string, 0, len(r))
	for n := range r {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintf(w, "  %-8s %s\n", n, r[n].Usage())
	}
}

// newFlagSet returns a flag set that already carries the shared --json flag.
func newFlagSet(name string) (*flag.FlagSet, *bool) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	json := fs.Bool("json", false, "machine-readable JSON output")
	return fs, json
}

// parseArgs parses fs while allowing flags to appear before or after positional
// arguments (Go's flag package otherwise stops at the first positional). It
// returns the positional arguments. Value-vs-bool arity is derived from the
// flag set, so there is no second list of flag names to keep in sync.
func parseArgs(fs *flag.FlagSet, args []string) ([]string, error) {
	bools := boolFlagNames(fs)
	var flagArgs, pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			pos = append(pos, args[i+1:]...)
			break
		}
		if len(a) < 2 || a[0] != '-' {
			pos = append(pos, a)
			continue
		}
		flagArgs = append(flagArgs, a)
		name := strings.TrimLeft(a, "-")
		if strings.IndexByte(name, '=') >= 0 || bools[name] {
			continue // inline value, or a boolean flag that takes none
		}
		if i+1 < len(args) {
			i++
			flagArgs = append(flagArgs, args[i])
		}
	}
	if err := fs.Parse(flagArgs); err != nil {
		return nil, err
	}
	return pos, nil
}

// boolFlagNames returns the names of the boolean flags in fs.
func boolFlagNames(fs *flag.FlagSet) map[string]bool {
	m := make(map[string]bool)
	fs.VisitAll(func(f *flag.Flag) {
		if bf, ok := f.Value.(interface{ IsBoolFlag() bool }); ok && bf.IsBoolFlag() {
			m[f.Name] = true
		}
	})
	return m
}

// multiFlag collects a repeatable string flag (e.g. --tag).
type multiFlag []string

func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error {
	*m = append(*m, v)
	return nil
}

// report prints either a JSON result or a human message.
func (cx *Context) report(ids []string, human string) error {
	if cx.JSON {
		return emitJSON(cx.Out, true, ids, nil, nil)
	}
	if human != "" {
		fmt.Fprintln(cx.Out, human)
	}
	return nil
}

// resolveSection maps a user-supplied section name (key, title, or a dashed
// variant) to a section key.
func resolveSection(cx *Context, name string) (string, bool) {
	low := strings.ToLower(strings.ReplaceAll(name, "-", "_"))
	for _, sec := range cx.Svc.Schema().Sections {
		if sec.Key == name || strings.ToLower(sec.Key) == low || strings.EqualFold(sec.Title, name) {
			return sec.Key, true
		}
	}
	return "", false
}

// candidates returns the display IDs an ID-or-query argument could refer to: a
// single exact ID when it resolves, otherwise every item whose task contains the
// query (case-insensitive).
func candidates(cx *Context, arg string) ([]string, error) {
	l, err := cx.Svc.List()
	if err != nil {
		return nil, err
	}
	s := cx.Svc.Schema()
	if _, err := l.Resolve(s, arg); err == nil {
		return []string{strings.ToUpper(strings.TrimSpace(arg))}, nil
	}
	q := strings.ToLower(arg)
	var matches []string
	for i := range l.Items {
		if id := l.ComputeID(s, i); id != "" && strings.Contains(strings.ToLower(l.Items[i].Task), q) {
			matches = append(matches, id)
		}
	}
	return matches, nil
}

// resolveArg turns an ID or fuzzy task query into a single item ID, erroring on
// no match or ambiguity.
func resolveArg(cx *Context, arg string) (string, error) {
	m, err := candidates(cx, arg)
	if err != nil {
		return "", err
	}
	switch len(m) {
	case 1:
		return m[0], nil
	case 0:
		return "", fmt.Errorf("no item matches %q", arg)
	default:
		return "", fmt.Errorf("%q is ambiguous: %s", arg, strings.Join(m, ", "))
	}
}
