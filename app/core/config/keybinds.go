package config

import (
	"fmt"
	"strings"
)

type KeyBindings map[string]KeyCommand

func (k KeyBindings) Validate() error {
	reserved := []string{
		"ctrl-c",
		"ctrl-m",
	}

	for key, cmd := range k {
		err := cmd.IsValid()
		if err != nil {
			return fmt.Errorf("invalid command for key %s: %w", key, err)
		}

		for _, r := range reserved {
			if key == r {
				return fmt.Errorf("key '%s' is reserved", key)
			}
		}
	}

	return nil
}

type KeyCommand struct {
	Cmd         string `toml:"cmd"`
	Desc string `toml:"desc"`
}

func (k KeyCommand) String() string {
	return string(k.Cmd)
}

func (k KeyCommand) IsValid() error {
	if strings.HasPrefix(":", k.Cmd) {
		validoptions := []string{":GitClone", ":GitPull", ":Exit"}

		// check if it's a valid option
		var found bool
		for _, option := range validoptions {
			if option == k.Cmd {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("invalid command '%s'", k.Cmd)
		}

		return nil
	}

	// assume that it's a shell command
	return nil
}
