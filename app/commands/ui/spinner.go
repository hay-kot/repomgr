package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	msg string
	ch  chan string

	spinner  spinner.Model
	quitting bool
}

var _ tea.Model = &Spinner{}

func NewSpinner(ch chan string, msg string) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

	spin := &Spinner{
		spinner: s,
		msg:     msg,
		ch:      ch,
	}

	go func() {
		for msg := range ch {
			spin.msg = msg
		}
	}()

	return spin
}

// Init implements tea.Model.
func (s *Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

// Update implements tea.Model.
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			s.quitting = true
			return s, tea.Quit
		default:
			return s, nil
		}

	default:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}
}

// View implements tea.Model.
func (s *Spinner) View() string {
	str := s.spinner.View() + " " + s.msg
	if s.quitting {
		return str + "\n"
	}
	return str
}
