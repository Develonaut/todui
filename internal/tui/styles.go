package tui

import "charm.land/lipgloss/v2"

// Palette. Numeric ANSI 256 colors keep todui readable on most terminals.
var (
	styleTitle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	styleDim      = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	styleID       = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	styleTag      = lipgloss.NewStyle().Foreground(lipgloss.Color("36"))
	styleClaim    = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	styleSelected = lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Bold(true)
	styleCursor   = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true)
	styleKey      = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	styleErr      = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	styleStatus   = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
	styleConfirm  = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
)

// sectionAccents cycles per section for the headers.
var sectionAccents = []string{"39", "78", "214", "244", "108"}

// sectionStyle returns a bold style for the section at display index i.
func sectionStyle(i int) lipgloss.Style {
	color := sectionAccents[i%len(sectionAccents)]
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))
}
