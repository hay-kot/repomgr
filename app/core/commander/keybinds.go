package commander

import (
	"fmt"
	"slices"
	"strings"
)

const (
	CmdStrOpen     = "open '{{ .Repo.HTMLURL }}'"
	CmdStrGitClone = "git clone '{{ .Repo.CloneSSHURL }}' '{{ .CloneDir }}'"
	CmdStrExit     = ":Exit '{{ .CloneDir }}'"
)

func NewDefaultKeyBindings() KeyBindings {
	return KeyBindings{
		"ctrl+o": KeyCommand{
			Cmd:  CmdStrOpen,
			Desc: "open url",
			Mode: ModeBackground,
		},
		"ctrl+p": KeyCommand{
			Cmd:  CmdStrGitClone,
			Desc: "clone repo",
			Mode: ModeReadOnly,
		},
		"enter": KeyCommand{
			Cmd:  CmdStrExit,
			Desc: "exit with clone directory path",
			Mode: ModeBackground,
		},
	}
}

type KeyBindings map[string]KeyCommand

func (k KeyBindings) Validate() error {
	reserved := []string{
		"ctrl+c",
		"ctrl+m",
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

type Mode string

const (
	ModeBackground  Mode = "background"
	ModeReadOnly    Mode = "readonly"
	ModeInteractive Mode = "interactive"
)

// KeyCommand defines a key binding command that can be executed
// by the Commander
type KeyCommand struct {
	Cmd  string `toml:"cmd"`
	Desc string `toml:"desc"`
	Mode Mode   `toml:"mode"`
}

func (k KeyCommand) String() string {
	return string(k.Cmd)
}

func (k KeyCommand) IsValid() error {
	if strings.HasPrefix(k.Cmd, ":") {
		validoptions := []string{
			":Exit",
		}

		found := slices.Contains(validoptions, k.Cmd)
		if !found {
			return fmt.Errorf("invalid command '%s'", k.Cmd)
		}

		return nil
	}

	// Validate mode
	modes := []Mode{
		ModeBackground,
		ModeReadOnly,
		ModeInteractive,
	}

	var found bool
	for _, mode := range modes {
		if mode == k.Mode {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid mode '%s'", k.Mode)
	}

	// assume that it's a shell command
	return nil
}
