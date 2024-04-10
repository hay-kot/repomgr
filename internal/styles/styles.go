package styles

import "github.com/charmbracelet/lipgloss"

const (
	ColorPrimary = "#3b82f6"
	ColorSuccess = "#22c55e"
	ColorError   = "#ef4444"
)

var (
	Bold    = lipgloss.NewStyle().Bold(true)
	Padding = lipgloss.NewStyle().PaddingLeft(1)
	Error   = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).PaddingLeft(1)
	Success = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).PaddingLeft(1)
)
