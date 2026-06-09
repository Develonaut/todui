// Package markdown imports a legacy Markdown task list into a todo.List. It is
// a best-effort, one-shot migration: line-oriented, lossy only where the source
// format is ambiguous, and faithful enough to round-trip through the Markdown
// renderer.
package markdown

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/develonaut/todui/internal/todo"
)

const claimedToken = "**CLAIMED**"

var (
	reID       = regexp.MustCompile(`^\*\*[A-Z]{1,3}\*\*\s+`)
	reLeadADO  = regexp.MustCompile(`^(#\d+)\s+`)
	reTrailTag = regexp.MustCompile("\\s*`\\[([^\\]]+)\\]`$")
	reDone     = regexp.MustCompile(`\s*\(done ([^)]*)\)\s*$`)
	reLastUpd  = regexp.MustCompile(`^_Last updated: (.*)_$`)
)

// Import parses src into a todo.List, mapping "## <title>" headings to schema
// sections by title. The list is normalized before return.
func Import(src []byte, schema todo.Schema) (todo.List, error) {
	titleToKey := make(map[string]string, len(schema.Sections))
	for _, sec := range schema.Sections {
		titleToKey[sec.Title] = sec.Key
	}

	var list todo.List
	order := make(map[string]int)
	current := ""
	inSections := false

	for _, line := range strings.Split(string(src), "\n") {
		switch {
		case strings.HasPrefix(line, "## "):
			inSections = true
			title := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			key, ok := titleToKey[title]
			if !ok {
				return todo.List{}, fmt.Errorf("importer: unknown section heading %q", title)
			}
			current = key
		case !inSections:
			if m := reLastUpd.FindStringSubmatch(line); m != nil {
				list.LastUpdated = m[1]
				continue
			}
			list.Header = append(list.Header, line)
		case strings.HasPrefix(line, "- [x] "):
			list.Items = append(list.Items, parseDone(strings.TrimPrefix(line, "- [x] "), current, order))
		case strings.HasPrefix(line, "- [ ] "):
			list.Items = append(list.Items, parseOpen(strings.TrimPrefix(line, "- [ ] "), current, order))
		}
	}

	trimTrailingBlanks(&list)
	list.Normalize(schema)
	return list, nil
}

// parseOpen extracts an open item: leading ID, optional CLAIMED marker, optional
// leading reference, and trailing tags. The remainder is the task (no attempt is
// made to split out context — that would be ambiguous).
func parseOpen(s, section string, order map[string]int) todo.Item {
	it := todo.Item{Section: section}
	s = reID.ReplaceAllString(s, "")
	if strings.HasPrefix(s, claimedToken) {
		it.Claimed = true
		s = strings.TrimSpace(strings.TrimPrefix(s, claimedToken))
	}
	if m := reLeadADO.FindStringSubmatch(s); m != nil {
		it.ADO = m[1]
		s = s[len(m[0]):]
	}
	for {
		m := reTrailTag.FindStringSubmatchIndex(s)
		if m == nil {
			break
		}
		it.Tags = append([]string{s[m[2]:m[3]]}, it.Tags...)
		s = s[:m[0]]
	}
	it.Title, it.Description = todo.SplitTitle(strings.TrimSpace(s))
	it.Order = order[section]
	order[section]++
	return it
}

// parseDone extracts a completed item and its trailing "(done ...)" annotation.
func parseDone(s, section string, order map[string]int) todo.Item {
	it := todo.Item{Section: section}
	if m := reDone.FindStringSubmatchIndex(s); m != nil {
		it.DoneDate = s[m[2]:m[3]]
		s = s[:m[0]]
	}
	it.Title, it.Description = todo.SplitTitle(strings.TrimSpace(s))
	it.Order = order[section]
	order[section]++
	return it
}

// trimTrailingBlanks drops blank lines at the end of the header; the renderer
// regenerates the separator before the first section.
func trimTrailingBlanks(list *todo.List) {
	for len(list.Header) > 0 && strings.TrimSpace(list.Header[len(list.Header)-1]) == "" {
		list.Header = list.Header[:len(list.Header)-1]
	}
}
