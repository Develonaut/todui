package cli

import (
	"fmt"
	"os"
	"time"

	mdimport "github.com/develonaut/todui/internal/importer/markdown"
)

type importCmd struct{}

func (importCmd) Name() string { return "import" }
func (importCmd) Usage() string {
	return "import --from FILE [--force]  — migrate a Markdown list into the store"
}

func (importCmd) Run(cx *Context, args []string) error {
	fs, jsonFlag := newFlagSet("import")
	from := fs.String("from", "", "source Markdown file")
	force := fs.Bool("force", false, "overwrite an existing store")
	if _, err := parseArgs(fs, args); err != nil {
		return err
	}
	cx.JSON = *jsonFlag

	if *from == "" {
		return fmt.Errorf("import: --from FILE required")
	}
	src, err := os.ReadFile(*from)
	if err != nil {
		return err
	}
	if !*force {
		if _, err := os.Stat(cx.Cfg.Store); err == nil {
			return fmt.Errorf("import: store %s already exists; pass --force to overwrite", cx.Cfg.Store)
		}
	}

	list, err := mdimport.Import(src, cx.Svc.Schema())
	if err != nil {
		return err
	}

	// Back up the source before any write: when the mirror points at the source
	// file, the first write will regenerate it.
	bak, err := backupFile(*from)
	if err != nil {
		return err
	}
	if err := cx.Svc.Replace(list); err != nil {
		return err
	}
	return cx.report(nil, fmt.Sprintf("Imported %d items from %s (backup: %s)", len(list.Items), *from, bak))
}

// backupFile copies path to a timestamped sibling and returns the backup path.
func backupFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	bak := fmt.Sprintf("%s.bak-%d", path, time.Now().Unix())
	if err := os.WriteFile(bak, b, 0o644); err != nil {
		return "", err
	}
	return bak, nil
}
