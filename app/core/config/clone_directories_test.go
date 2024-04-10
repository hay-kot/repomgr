package config

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func Test_CloneDirectories_Validate(t *testing.T) {
	type tcase struct {
		name    string
		input   CloneDirectories
		wantErr bool
	}

	cases := []tcase{
		{
			name: "default is required",
			input: CloneDirectories{
				Default:  "",
				Matchers: []Matcher{},
			},
			wantErr: true,
		},
		{
			name: "blank matcher glob",
			input: CloneDirectories{
				Default: "clone-dir",
				Matchers: []Matcher{
					{Match: "", Directory: ""},
				},
			},
			wantErr: true,
		},
		{
			name: " invalid matcher glob",
			input: CloneDirectories{
				Default: "clone-dir",
				Matchers: []Matcher{
					{Match: "[", Directory: "dir"},
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			err := tc.input.Validate()

			if tc.wantErr {
				is.True(err != nil)
			} else {
				is.NoErr(err)
			}
		})
	}
}

func Test_Matcher_IsMatch(t *testing.T) {
	type tcase struct {
		match string
		input string
		want  bool
	}

	cases := []tcase{
		{match: "foo", input: "foo", want: true},
		{match: "foo/*", input: "foo/bar", want: true},
		{match: "*/*", input: "foo/bar", want: true},
		{match: "*", input: "anything", want: true},
		{match: "exact", input: "exacts", want: false},
	}

	for _, tc := range cases {
		is := is.New(t)
		t.Run(fmt.Sprintf("match=%s input=%s", tc.match, tc.input), func(t *testing.T) {
			is := is.New(t)

			m := Matcher{Match: tc.match}
			got := m.IsMatch(tc.input)

			is.Equal(got, tc.want)
		})
	}
}
