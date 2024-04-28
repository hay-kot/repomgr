package commander

import (
	"errors"
	"io"
	"os/exec"
	"testing"

	"github.com/matryer/is"
)

type tActionCommand struct {
	run func() error
	set func(w io.Writer)
}

func (tac *tActionCommand) Run() error {
	return tac.run()
}

func (tac *tActionCommand) SetWriter(w io.Writer) {
	tac.set(w)
}

func Test_Action_Run(t *testing.T) {
	tests := []struct {
		name           string
		runFunc        func() error
		expectedErrMsg string
	}{
		{
			name: "WithError",
			runFunc: func() error {
				return errors.New("WithError error msg")
			},
			expectedErrMsg: "WithError error msg",
		},
		{
			name: "WithoutError",
			runFunc: func() error {
				return nil
			},
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			action := &Action{
				Mode: ModeBackground,
				cmd: &tActionCommand{
					run: tt.runFunc,
					set: func(w io.Writer) {},
				},
			}

			err := <-action.GoRun()

			is := is.New(t)
			if tt.expectedErrMsg != "" {
				is.True(err != nil)                      // want error to be not nil
				is.Equal(err.Error(), tt.expectedErrMsg) // want specific error message
			} else {
				is.True(err == nil) // want error to be nil
			}
		})
	}
}

func Test_Action_IsExec(t *testing.T) {
	tests := []struct {
		name       string
		action     *Action
		expectedOk bool
	}{
		{
			name: "is exec command",
			action: &Action{
				Mode: ModeBackground,
				cmd: execWrapper{
					Cmd: &exec.Cmd{},
				},
			},
			expectedOk: true,
		},
		{
			name: "is not exec command",
			action: &Action{
				Mode: ModeBackground,    // Different mode
				cmd:  &tActionCommand{}, // No command
			},
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cmd, ok := tt.action.IsExec()

			is := is.New(t)
			is.Equal(ok, tt.expectedOk)         // Check if ok matches expectedOk
			is.Equal(cmd != nil, tt.expectedOk) // Check if cmd is not nil when expectedOk is true
		})
	}
}
