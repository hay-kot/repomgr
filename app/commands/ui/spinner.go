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

// NewSpinnerFunc creates a new spiner with a given initial message and a function that
// will be executed in the background in a goroutine. The function can use the provided
// channel to send message stop the spinner to update the out. When the function returns
// the spinner will stop and the function will return the error result from the function.
func NewSpinnerFunc(initial string, fn func(ch chan<- string) error) error {
	var (
		done  = make(chan struct{})
		msgch = make(chan string)
	)

	spinner := NewSpinner(msgch, done, initial)

	var outerr error
	go func() {
		outerr = fn(msgch)
		close(done)
	}()

	p := tea.NewProgram(spinner)
	_, err := p.Run()
	if err != nil {
		return err
	}

	return outerr
}

func NewSpinner(ch chan string, done chan struct{}, msg string) *Spinner {
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

	go func() {
		<-done
		spin.quitting = true
	}()

	return spin
}

// Init implements tea.Model.
func (s *Spinner) Init() tea.Cmd {
	return s.spinner.Tick
}

// Update implements tea.Model.
func (s *Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s.quitting {
		return s, tea.Quit
	}

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
