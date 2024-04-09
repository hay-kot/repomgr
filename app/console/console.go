package console

import (
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	stylePadding = lipgloss.NewStyle().PaddingLeft(1)
	styleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("#b91c1c")).PaddingLeft(1)
	styleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("#14532d")).PaddingLeft(1)
)

type Console struct {
	writer io.Writer
	color  bool
}

func NewConsole(writer io.Writer, color bool) *Console {
	return &Console{
		writer: writer,
		color:  color,
	}
}

func (c *Console) write(str string) {
	_, err := c.writer.Write([]byte(str))
	if err != nil {
		panic(err)
	}
}

func (c *Console) UnknownError(title string, err error) {
	bldr := strings.Builder{}

	bldr.WriteString(styleError.Render("An unexpected error occurred"))
	bldr.WriteString("\n\n")
	bldr.WriteString(stylePadding.Render("error"))
	bldr.WriteString("\n  '")
	bldr.WriteString(err.Error())
	bldr.WriteString("'\n")

	c.write(bldr.String())
}
