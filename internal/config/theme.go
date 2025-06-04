package config

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	PromptColor      lipgloss.Style
	InputColor       lipgloss.Style
	CorrectColor     lipgloss.Style
	IncorrectColor   lipgloss.Style
	ProgressColor    lipgloss.Style
	QuitMessageColor lipgloss.Style
}
