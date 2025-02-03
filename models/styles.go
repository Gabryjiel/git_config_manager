package models

import "github.com/charmbracelet/lipgloss"

var (
	headerStyle      = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
	arrowCinStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(14))
	localValueStyle  = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
	globalValueStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(11))
	systemValueStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(12))
	CenterStyle      = lipgloss.NewStyle().Width(80).AlignHorizontal(lipgloss.Center)
)
