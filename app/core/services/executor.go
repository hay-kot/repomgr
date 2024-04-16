package services

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

type Executor interface {
	Execute(cmd string, args ...string) error
	ExecuteHandler(cmd string, args ...string) (*CommandHandle, error)
}

type ShellExecutor struct {
	shell string
}

func NewShellExecutor(shell string) ShellExecutor {
	return ShellExecutor{shell: shell}
}

func (e ShellExecutor) Execute(cmd string, args ...string) error {
	log.Debug().Str("cmd", cmd).Strs("args", args).Msg("executing command")
	err := exec.Command(e.shell, "-c", cmd+" "+strings.Join(args, " ")).Run()
	log.Debug().Err(err).Msg("command executed")
	return err
}

func (e ShellExecutor) ExecuteHandler(cmd string, args ...string) (*CommandHandle, error) {
	c := exec.Command(e.shell, "-c", cmd+" "+strings.Join(args, " "))

	cw := &CommandHandle{cmd: c}

	c.Stderr = cw
	c.Stdout = cw

	return cw, nil
}

type CommandHandle struct {
	cmd  *exec.Cmd
	bits [][]byte
}

func (cw *CommandHandle) Run() error {
	return cw.cmd.Run()
}

func (cw *CommandHandle) Start() error {
	return cw.cmd.Start()
}

func (cw *CommandHandle) Wait() error {
	return cw.cmd.Wait()
}

// Write writes the output of the command to the channel
func (cw *CommandHandle) Write(p []byte) (n int, err error) {
	log.Debug().Str("output", string(p)).Msg("command output")
	cw.bits = append(cw.bits, p)
	return len(p), nil
}

func (cw *CommandHandle) String() string {
	return string(bytes.Join(cw.bits, []byte{}))
}
