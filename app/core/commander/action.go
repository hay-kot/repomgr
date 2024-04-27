package commander

import (
	"io"
	"os/exec"
)

type Action struct {
	Mode        Mode
	cmd         ActionCommand
	onFinished  []func()
	isExit      bool
	exitMessage string
}

// IsExec returns true if the action is an exec action. It will also
// return the underlying exec.Cmd object.
func (a *Action) IsExec() (*exec.Cmd, bool) {
	v, ok := a.cmd.(execWrapper)
	if !ok {
		return nil, false
	}

	return v.Cmd, true
}

// IsExit returns true if the action is an exit action.
// If true, the ExitMessage will contain the message to display.
func (a *Action) IsExit() bool {
	return a.isExit
}

// ExitMessage returns the message to display when the action is an exit action.
func (a *Action) ExitMessage() string {
	return a.exitMessage
}

func (a *Action) SetWriter(w io.Writer) *Action {
	if a.cmd != nil {
		a.cmd.SetWriter(w)
	}

	return a
}

func (s *Action) OnFinished(fn func()) *Action {
	s.onFinished = append(s.onFinished, fn)
	return s
}

// GoRun runs the action within a go-routine and returns an err-channel
// that will be closed when the action is finished with the resulting error
// if any.
func (a *Action) GoRun() <-chan error {
	errch := make(chan error, 1)

	go func() {
		err := a.cmd.Run()
		if err != nil {
			errch <- err
		}

		for _, fn := range a.onFinished {
			fn()
		}

		close(errch)
	}()

	return errch
}
