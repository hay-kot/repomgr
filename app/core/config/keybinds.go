package config

import (
	"fmt"
	"strings"
)

type KeyBindings map[string]KeyCommand

func (k KeyBindings) Validate() error {
	for key, cmd := range k {
		err := cmd.IsValid()
		if err != nil {
			return fmt.Errorf("invalid command for key %s: %w", key, err)
		}
	}

	return nil
}

type KeyCommand string

func (k KeyCommand) String() string {
	return string(k)
}

func (k KeyCommand) IsValid() error {
	str := string(k)
	if strings.HasPrefix("::", str) {
		validoptions := []string{"::open", "::clone", "::shell"}

		// check if it's a valid option
		var found bool
		for _, option := range validoptions {
			if option == str {
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("invalid command '%s'", str)
		}

		return nil
	}

	// assume that it's a shell command
	return nil
}
