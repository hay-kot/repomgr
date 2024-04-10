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

const (
	ColorBlue   = lipgloss.Color("#255F85")
	ColorRed    = lipgloss.Color("#DA4167")
	ColorSubtle = lipgloss.Color("#848484")
	ColorWhite  = lipgloss.Color("#FFFFFF")
)

var (
	Subtle       = lipgloss.NewStyle().Foreground(ColorSubtle).Render
	AccentRed    = lipgloss.NewStyle().Foreground(ColorRed).Render
	AccentBlue   = lipgloss.NewStyle().Foreground(ColorBlue).Render
	HighlightRow = lipgloss.NewStyle().Background(lipgloss.Color("#2D2F27")).Render
)
