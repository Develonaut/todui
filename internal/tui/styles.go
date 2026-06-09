package tui

import "charm.land/lipgloss/v2"

// A vibrant, Charm/Crush-flavored palette: magenta and purple accents on a dark
// background, with soft-lavender body text. The detail body stays dimmer than
// the list so the selected row remains the focus.
var (
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C6C7E")) // metadata, counts
	styleFaint   = lipgloss.NewStyle().Foreground(lipgloss.Color("#44444E")) // separators
	styleID      = lipgloss.NewStyle().Foreground(lipgloss.Color("#9D7CFF")) // ids, violet
	styleItem    = lipgloss.NewStyle().Foreground(lipgloss.Color("#C8C8D8")) // unselected titles
	styleSelect  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA")).Bold(true)
	styleCursor  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FD2")) // cursor arrow, magenta
	styleTag     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D7AF")) // tags, teal
	styleDetail  = lipgloss.NewStyle().Foreground(lipgloss.Color("#9A9AB0")) // detail body
	styleKey     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FD2")) // help keys, magenta
	styleErr     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F87"))
	styleStatus  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00D787"))
	styleConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD75F"))

	styleBorderActive = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FD2")) // focused frame, magenta
	styleBorder       = lipgloss.NewStyle().Foreground(lipgloss.Color("#44444E")) // inactive frame
	styleLabel        = lipgloss.NewStyle().Foreground(lipgloss.Color("#9D7CFF")).Bold(true)
)

// sectionAccents give each section a vibrant header color (the emoji carry the
// cue too). Order matches the typical In Progress / Now / Next / Later / Done.
var sectionAccents = []string{"#FF5FD2", "#00D787", "#00E5FF", "#9D7CFF", "#6C6C7E"}

// sectionStyle returns the header color for the section at display index i.
func sectionStyle(i int) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(sectionAccents[i%len(sectionAccents)]))
}

// logoColors are the per-letter gradient stops for the TODUI wordmark
// (magenta → purple).
var logoColors = []string{"#FF06B7", "#E63DD9", "#C84BEC", "#A453F4", "#7D56F4"}
