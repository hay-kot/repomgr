package ui

import (
	"bytes"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/repomgr/app/core/commander"
	"github.com/hay-kot/repomgr/internal/icons"
	"github.com/hay-kot/repomgr/internal/styles"
)

type WriteEvent struct{}

type ListenableWriter struct {
	bytes  [][]byte
	notify chan struct{}
}

func (w *ListenableWriter) Write(p []byte) (n int, err error) {
	w.bytes = append(w.bytes, p)

	select {
	case w.notify <- struct{}{}:
	default:
	}

	return len(p), nil
}

func (w *ListenableWriter) WaitUntilContent() WriteEvent {
	<-w.notify
	return WriteEvent{}
}

func (w *ListenableWriter) String() string {
	return string(bytes.Join(w.bytes, []byte{}))
}

func (w *ListenableWriter) Lines() []string {
	out := make([]string, len(w.bytes))
	for i, b := range w.bytes {
		out[i] = string(b)
	}

	return out
}

var _ tea.Model = &CommandView{}

type CommandView struct {
	spinner  spinner.Model
	init     bool
	finished bool
	handle   *commander.Action
	back     tea.Model
	errch    <-chan error
	writer   *ListenableWriter
	donech   chan struct{}
}

type CommandFinishedMsg struct{}

func NewCommandView(handle *commander.Action, back tea.Model) *CommandView {
	writer := &ListenableWriter{notify: make(chan struct{})}
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.ColorPrimary))

	c := &CommandView{
		spinner: s,
		writer:  writer,
		back:    back,
		donech:  make(chan struct{}),
	}

	c.handle = handle.
		SetWriter(writer).
		OnFinished(func() {
			c.donech <- struct{}{}
		})

	return c
}

// Init implements tea.Model.
func (c *CommandView) Init() tea.Cmd {
	c.init = true
	c.errch = c.handle.GoRun()

	return tea.Batch(
		c.spinner.Tick,
		func() tea.Msg {
			return c.writer.WaitUntilContent()
		},
		func() tea.Msg {
			<-c.donech
			c.finished = true
			return CommandFinishedMsg{}
		},
	)
}

// Update implements tea.Model.
func (c *CommandView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := tea.Batch(func() tea.Msg {
		return c.writer.WaitUntilContent()
	})

	if !c.init {
		cmd := c.Init()
		cmds = tea.Batch(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if c.finished {
				return c.back, nil
			}
		case "q", "esc", "ctrl+c":
			return c.back, nil
		}
	}

	var spincmd tea.Cmd
	c.spinner, spincmd = c.spinner.Update(msg)
	return c, tea.Batch(cmds, spincmd)
}

// View implements tea.Model.
func (c *CommandView) View() string {
	bldr := &strings.Builder{}

	const (
		Offset    = "  "
		CmdOffset = Offset + "  > "
	)

	if !c.finished {
		bldr.WriteString(Offset)
		bldr.WriteString(c.spinner.View() + " " + "Executing")
	} else {
		bldr.WriteString(" ") // Not using offset because of double width characters
		bldr.WriteString(styles.Success.Render(icons.Check + " "))
		bldr.WriteString("Success")
	}
	bldr.WriteString("\n")

	for _, line := range c.writer.Lines() {
		line = strings.TrimSpace(line)
		bldr.WriteString(CmdOffset)
		bldr.WriteString(styles.Subtle(line))

		if line != "" {
			bldr.WriteString("\n")
		}
	}

	if c.finished {
		bldr.WriteString("\n")
		bldr.WriteString(Offset)
		bldr.WriteString("Press Enter to continue")
	}

	return bldr.String()
}
