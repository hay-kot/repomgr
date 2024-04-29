package commander

import (
	"io"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type ActionCommandBuilder interface {
	Build(cmd string, args ...string) ActionCommand
}

type ActionCommand interface {
	Run() error
	SetWriter(w io.Writer)
}

var _ ActionCommandBuilder = &ShellCommandBuilder{}

// ShellCommandBulder implements the ActionCommandBuilder interface for shell
// commands. It returns a new ActionCommand that will execute the given command
// with the given arguments in a shell using the exec package.
type ShellCommandBuilder struct {
	Shell        string
	ShellCmdFlag string
}

// Build implements ActionCommandBuilder.
func (s *ShellCommandBuilder) Build(cmd string, args ...string) ActionCommand {
	log.Debug().
		Str("shell", s.Shell).
		Str("cmd", cmd).
		Strs("args", args).
		Msg("executing command")

	c := exec.Command(s.Shell, s.ShellCmdFlag, cmd+" "+strings.Join(args, " "))
	return execWrapper{Cmd: c}
}
