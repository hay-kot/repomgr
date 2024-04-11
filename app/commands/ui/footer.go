package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/core/config"
)

var _ tea.Model = &Footer{}

type Footer struct {
	keys config.KeyBindings
}

func NewFooter(keybindings config.KeyBindings) *Footer {
	return &Footer{
		keys: keybindings,
	}
}

// Init implements tea.Model.
func (f *Footer) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (f *Footer) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return f, nil
}

// View implements tea.Model.
func (f *Footer) View() string {
	return "Hello World"
}
