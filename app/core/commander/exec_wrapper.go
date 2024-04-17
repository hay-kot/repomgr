package commander

import (
	"io"
	"os/exec"
)

var _ ActionCommand = &execWrapper{}

// execWrapper wraps an exec.Cmd to implement the ActionCommand interface
type execWrapper struct {
	*exec.Cmd
}

func (e execWrapper) SetWriter(w io.Writer) {
	e.Stdout = w
	e.Stderr = w
}
