package tui

import "charm.land/lipgloss/v2"

// A deliberately quiet, muted palette for a dark terminal. Body text sits in
// soft grays; accents are desaturated so nothing shouts. The detail pane is
// dimmer than the list so the selected row stays the focus.
var (
	styleTitle   = lipgloss.NewStyle().Foreground(lipgloss.Color("110")) // app name, soft blue
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // metadata, counts
	styleFaint   = lipgloss.NewStyle().Foreground(lipgloss.Color("237")) // rules, separators
	styleID      = lipgloss.NewStyle().Foreground(lipgloss.Color("67"))  // ids, muted blue
	styleItem    = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // unselected titles, soft gray
	styleSelect  = lipgloss.NewStyle().Foreground(lipgloss.Color("253")) // selected title, soft white
	styleCursor  = lipgloss.NewStyle().Foreground(lipgloss.Color("110")) // cursor arrow
	styleClaim   = lipgloss.NewStyle().Foreground(lipgloss.Color("137")) // claimed dot, muted amber
	styleTag     = lipgloss.NewStyle().Foreground(lipgloss.Color("66"))  // tags, muted teal
	styleDetail  = lipgloss.NewStyle().Foreground(lipgloss.Color("244")) // detail body, subtle
	styleKey     = lipgloss.NewStyle().Foreground(lipgloss.Color("245")) // help keys
	styleErr     = lipgloss.NewStyle().Foreground(lipgloss.Color("174")) // muted red
	styleStatus  = lipgloss.NewStyle().Foreground(lipgloss.Color("108")) // muted green
	styleConfirm = lipgloss.NewStyle().Foreground(lipgloss.Color("179")) // muted amber

	styleBorder       = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // inactive frame
	styleBorderActive = lipgloss.NewStyle().Foreground(lipgloss.Color("67"))  // focused frame
	styleLabel        = lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Bold(true)
)

// sectionAccents desaturate the section headers (the emoji already carry the
// strong color cue, so the text stays quiet).
var sectionAccents = []string{"67", "108", "137", "244", "108"}

// sectionStyle returns the muted header style for the section at display index i.
func sectionStyle(i int) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(sectionAccents[i%len(sectionAccents)]))
}
