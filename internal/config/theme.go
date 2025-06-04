package config

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	PromptColor    lipgloss.Style
	InputColor     lipgloss.Style
	CorrectColor   lipgloss.Style
	IncorrectColor lipgloss.Style
	ProgressColor  lipgloss.Style
	QuitColor      lipgloss.Style
	WelcomeColor   lipgloss.Style
}

func DefaultTheme() Theme {
	return Theme{
		PromptColor:    lipgloss.NewStyle().Foreground(lipgloss.Color("#5cf043")), // Orange
		InputColor:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")), // White
		CorrectColor:   lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")), // Green
		IncorrectColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4C4C")), // Red
		ProgressColor:  lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB")), // Sky Blue
		QuitColor:      lipgloss.NewStyle().Foreground(lipgloss.Color("#808080")), // Gray
		WelcomeColor:   lipgloss.NewStyle().Foreground(lipgloss.Color("#be6aff")), // Gold
	}
}
