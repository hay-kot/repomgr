package quicktmpl

import (
	"testing"

	"github.com/matryer/is"
)

func Test_New(t *testing.T) {
	type tcase struct {
		name    string
		tmpl    string
		wantErr bool
	}

	cases := []tcase{
		{
			name: "valid template",
			tmpl: "Hello, {{ .Name }}!",
		},
		{
			name:    "invalid template",
			tmpl:    "Hello, {{ .Name!",
			wantErr: true,
		},
	}

	is := is.New(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			t.Parallel()
			_, err := New(tc.tmpl)

			if tc.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
		})
	}
}

func Test_Render(t *testing.T) {
	type tcase struct {
		name    string
		tmpl    string
		data    Data
		want    string
		wantErr bool
	}

	cases := []tcase{
		{
			name: "valid template",
			tmpl: "Hello, {{ .Name }}!",
			data: Data{"Name": "World"},
			want: "Hello, World!",
		},
		{
			name:    "invalid template",
			tmpl:    "Hello, {{ .Name!",
			data:    Data{"Name": "World"},
			wantErr: true,
		},
		{
			name: " template with function",
			tmpl: "Hello, {{ .Name | trim |  lower }}!",
			data: Data{"Name": " World "},
			want: "Hello, world!",
		},
	}

	is := is.New(t)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			is := is.New(t)

			t.Parallel()
			got, err := Render(tc.tmpl, tc.data)

			if tc.wantErr {
				is.True(err != nil) // want error
				return
			}

			is.NoErr(err)
			is.Equal(got, tc.want)
		})
	}
}
