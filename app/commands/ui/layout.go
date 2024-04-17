package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Layout{}

type Layout struct {
	Body tea.Model
}

func NewLayout(body tea.Model) *Layout {
	return &Layout{
		Body: body,
	}
}

// Init implements tea.Model.
func (l *Layout) Init() tea.Cmd {
	return l.Body.Init()
}

// Update implements tea.Model.
func (l *Layout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return l.Body.Update(msg)
}

// View implements tea.Model.
func (l *Layout) View() string {
	return l.Body.View()
}
