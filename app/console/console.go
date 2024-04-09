package console

import (
	"io"
	"strings"

	"github.com/hay-kot/repomgr/internal/icons"
	"github.com/hay-kot/repomgr/internal/styles"
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

	bldr.WriteString(styles.Error.Render("An unexpected error occurred"))
	bldr.WriteString("\n\n")
	bldr.WriteString(styles.Padding.Render("Error"))
	bldr.WriteString("\n  '")
	bldr.WriteString(err.Error())
	bldr.WriteString("'\n")

	c.write(bldr.String())
}

type ListItem struct {
	StatusOk bool
	Status   string
}

func (c *Console) List(title string, items []ListItem) {
	bldr := strings.Builder{}

	bldr.WriteString(styles.Bold.Render(styles.Padding.Render(title)))
	bldr.WriteString("\n")

	for _, item := range items {
		bldr.WriteString(" ")
		if item.StatusOk {
			bldr.WriteString(
				styles.Success.Render(icons.Check),
			)
		} else {
			bldr.WriteString(styles.Error.Render(icons.Cross))
		}

		bldr.WriteString(" ")
		bldr.WriteString(item.Status)
		bldr.WriteString("\n")
	}

	c.write(bldr.String())
}

func (c *Console) LineBreak() {
	c.write("\n")
}
