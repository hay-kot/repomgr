package commander

import (
	"strings"
	"testing"

	"github.com/matryer/is"
)

func Test_KeyCommand_IsValid(t *testing.T) {
	type tcase struct {
		name        string
		cmd         KeyCommand
		wantErr     bool
		errContains string
	}

	cases := []tcase{
		{
			name: "valid command",
			cmd: KeyCommand{
				Cmd:  ":Exit {{ .Repo.CloneDir }}",
				Desc: "Exit the application",
				Mode: ModeBackground,
			},
		},
		{
			name: "invalid mode",
			cmd: KeyCommand{
				Cmd:  "git clone",
				Desc: "Exit the application",
				Mode: "test",
			},
			wantErr:     true,
			errContains: "mode",
		},
		{
			name: "invalid command",
			cmd: KeyCommand{
				Cmd:  ":NotValid",
				Desc: "Exit the application",
				Mode: ModeBackground,
			},
			wantErr:     true,
			errContains: "command",
		},
	}

	is := is.New(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			err := tc.cmd.IsValid()

			if tc.wantErr {
				is.True(err != nil)                                    // expected error
				is.True(strings.Contains(err.Error(), tc.errContains)) // error message should contain
				return
			}

			is.NoErr(err) // no error expected
		})
	}
}
