package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Title styles
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1)

	// Status styles
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("214"))

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	// Content styles
	Bold = lipgloss.NewStyle().
		Bold(true)

	Dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	Subtle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	// Special styles
	Command = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212"))

	Path = lipgloss.NewStyle().
		Foreground(lipgloss.Color("147")).
		Italic(true)

	Badge = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("42")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	// Prompt styles
	Prompt = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	// Code block style
	Code = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)
)

// Helper functions
func ErrorMsg(msg string) string {
	return Error.Render("✗ " + msg)
}

func SuccessMsg(msg string) string {
	return Success.Render("✓ " + msg)
}

func InfoMsg(msg string) string {
	return Info.Render("→ " + msg)
}

func WarningMsg(msg string) string {
	return Warning.Render("⚠ " + msg)
}