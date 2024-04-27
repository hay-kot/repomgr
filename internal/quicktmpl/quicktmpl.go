package quicktmpl

import (
	"bytes"
	"strings"
	"text/template"
)

// Data is a alias for a map of string to any and is used
// for passing data to the template.
type Data map[string]any

// Funcs is a map of functions that can be used in the template.
// you can modify this map to add more functions to the template.
// This should be done before calling New or Render.
var Funcs = template.FuncMap{
	"trim":    strings.TrimSpace,
	"trimset": strings.Trim,
	"join":    strings.Join,
	"lower":   strings.ToLower,
	"upper":   strings.ToUpper,
}

// New creates a new template with the given template string. It also adds the
// Funcs to the template.
func New(tmpl string) (*template.Template, error) {
	t, err := template.New("t").Funcs(Funcs).Parse(tmpl)

	return t, err
}

// Render renders a template with the given data. It returns the rendered
// template as a string. This is equivalent to calling New and Execute on the
// template.
func Render(tmpl string, data Data) (string, error) {
	t, err := New(tmpl)
	if err != nil {
		return "", err
	}

	b := bytes.NewBuffer(make([]byte, 0, len(tmpl)))
	err = t.Execute(b, data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
