package config

import (
	"os"
	"path/filepath"
)

func ExpandPath(confpath, input string) string {
	if input[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		input = home + input[1:]
	}

	if input[:2] == "./" {
		confdir := filepath.Dir(confpath)
		input = confdir + input[1:]
	}

	return input
}

type validator interface {
	Validate() error
}
