package services

import (
	"strings"
	"testing"

	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/matryer/is"
)

type executeRecorder struct {
	cmd  string
	args []string
}

func (e *executeRecorder) Execute(cmd string, args ...string) error {
	e.cmd = cmd
	e.args = args
	return nil
}

func (e *executeRecorder) reset() {
	e.cmd = ""
	e.args = nil
}

func (e *executeRecorder) String() string {
	return strings.TrimSpace(e.cmd + " " + strings.Join(e.args, " "))
}

func Test_CommandService_Run(t *testing.T) {
	repo := factory(1)[0]

	type tcase struct {
		name          string
		input         string
		want          string
		commandResult *CommandResult
	}

	e := &executeRecorder{}

	cfg := config.Default()

	cfg.CloneDirectories = config.CloneDirectories{
		Default:  "/tmp/{{ .Repo.Username }}/{{ .Repo.Name }}",
		Matchers: []config.Matcher{},
	}

	s := tAppService(t, tAppServiceOpts{
		recorder: e,
		cfg:      cfg,
	})

	tcases := []tcase{
		{
			name:  "basic command with template",
			input: "open {{ .Repo.HTMLURL }}",
			want:  "open " + repo.HTMLURL,
		},
		{
			name:  "app command :GitClone",
			input: "git clone '{{ .Repo.CloneSSHURL }}' '{{ .CloneDir }}'",
			want:  "git clone '" + repo.CloneSSHURL + "' '/tmp/" + repo.DisplayName() + "'",
		},
		{
			name:  "app command :Exit",
			input: ":Exit '{{ .CloneDir }}'",
			want:  "",
			commandResult: &CommandResult{
				IsExit:      true,
				ExitMessage: "'/tmp/" + repo.DisplayName() + "'",
			},
		},
	}

	is := is.New(t)
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			e.reset()
			cr, err := s.Run(repo, tc.input)
			is.NoErr(err) // expected no error on execute
			is.Equal(e.String(), tc.want)

			if tc.commandResult != nil {
				is.Equal(*tc.commandResult, cr)
			}
		})
	}
}
