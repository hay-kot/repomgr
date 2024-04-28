package config

import (
	"testing"

	"github.com/matryer/is"
)

func Test_Source_TokenLoader(t *testing.T) {
	type tcase struct {
		name   string
		source Source
		want   string
		hook   func(t *testing.T)
	}

	tcases := []tcase{
		{
			name:   "no prefix",
			source: Source{TokenKey: "token value"},
			want:   "token value",
			hook:   func(t *testing.T) {},
		},
		{
			name:   "env prefix",
			source: Source{TokenKey: "env:TOKEN"},
			want:   "TEST_TOKEN_VALUE",
			hook: func(t *testing.T) {
				t.Setenv("TOKEN", "TEST_TOKEN_VALUE")
			},
		},
		{
			name:   "empty env",
			source: Source{TokenKey: "env:TOKEN"},
			want:   "",
			hook:   func(t *testing.T) {},
		},
	}

	is := is.New(t)
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			tc.hook(t)
			is.Equal(tc.source.Token(), tc.want) // loaded token should match expected
		})
	}
}

func Test_Source_Validate(t *testing.T) {
	type tcase struct {
		name    string
		source  Source
		wantErr bool
	}

	tcases := []tcase{
		{
			name:    "empty source",
			source:  Source{},
			wantErr: true,
		},
		{
			name: "no username",
			source: Source{
				Type:     SourceTypeGithub,
				TokenKey: "token",
				Username: "",
			},
			wantErr: true,
		},
		{
			name: "invalid source type",
			source: Source{
				Type:     "invalid",
				TokenKey: "token",
				Username: "username",
			},
			wantErr: true,
		},
		{
			name: "valid source",
			source: Source{
				Type:     SourceTypeGithub,
				TokenKey: "token",
				Username: "username",
			},
			wantErr: false,
		},
	}

	is := is.New(t)

	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			err := tc.source.Validate()
			if tc.wantErr {
				is.True(err != nil)
			} else {
				is.NoErr(err)
			}
		})
	}
}
