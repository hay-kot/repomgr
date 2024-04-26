package commander

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func Test_AppAction_IsValid(t *testing.T) {
	type tcase struct {
		name   string
		action AppAction
		want   bool
	}

	cases := []tcase{
		{
			name:   "invalid action",
			action: AppAction("invalid"),
			want:   false,
		},
		{
			name:   "invalid action",
			action: AppAction(":InvalidOp"),
			want:   false,
		},
	}

	valid := []AppAction{AppActionFork, AppActionExit}
	for _, v := range valid {
		cases = append(cases, tcase{
			want:   true,
			name:   fmt.Sprintf("valid action (%s)", v),
			action: v,
		})
	}

	is := is.New(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(tc.action.IsValid(), tc.want)
		})
	}
}

func Test_ParseAppAction(t *testing.T) {
	type tcase struct {
		name       string
		input      string
		wantAction AppAction
		wantArg    string
		wantOk     bool
	}

	cases := []tcase{
		{
			name:   "empty input",
			input:  "",
			wantOk: false,
		},
		{
			name:   "invalid input",
			input:  "invalid",
			wantOk: false,
		},
		{
			name:       "valid input",
			input:      ":GitFork",
			wantAction: AppActionFork,
			wantOk:     true,
		},
		{
			name:       "valid input with args",
			input:      ":Exit some message",
			wantAction: AppActionExit,
			wantArg:    "some message",
			wantOk:     true,
		},
	}

	is := is.New(t)
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s (%s)", tc.name, tc.input), func(t *testing.T) {
			is := is.New(t)
			action, args, ok := ParseAppAction(tc.input)

			is.Equal(action, tc.wantAction)
			is.Equal(args, tc.wantArg)
			is.Equal(ok, tc.wantOk)
		})
	}
}
