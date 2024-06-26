package commander

import (
	"strings"

	"github.com/hay-kot/repomgr/app/core/repofs"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/hay-kot/repomgr/internal/quicktmpl"
	"github.com/rs/zerolog/log"
)

type Commander struct {
	bindings KeyBindings
	rfs      *repofs.RepoFS
	bldr     ActionCommandBuilder
}

func New(bindings KeyBindings, rfs *repofs.RepoFS, bldr ActionCommandBuilder) *Commander {
	return &Commander{
		bindings: bindings,
		rfs:      rfs,
		bldr:     bldr,
	}
}

func (c *Commander) Bindings() KeyBindings {
	return c.bindings
}

// GetAction returns an action for a given key binding. If the key binding is not found, it will
// attempt to render the command template and return an action with the rendered command.
//
// Special Cases
//
//   - ":Exit" - If the command starts with ":Exit", we return an exit action with the message that
//     follows the command.
func (c *Commander) GetAction(key string, repo repos.Repository) (action *Action, ok bool) {
	commandTmpl, ok := c.bindings[key]
	if !ok {
		log.Debug().Str("key", key).Msg("key not found in bindings")
		return nil, false
	}

	cmdRendered, err := c.renderCommandTemplate(repo, commandTmpl.Cmd)
	if err != nil {
		return nil, false
	}

	if MatchesActionSyntax(cmdRendered) {
		appAction, rest, ok := ParseAppAction(cmdRendered)
		if !ok {
			log.Error().Str("cmd", cmdRendered).Msg("invalid app action")
			return nil, false
		}

		switch appAction {
		case AppActionFork:
			// do something
			panic("not implemented")
		case AppActionExit:
			// Special case for `Exit: ...` if the command starts with ":Exit", we return an exit action
			// with the message that follows the command.
			return &Action{
				isExit:      true,
				exitMessage: rest,
			}, true
		}
	}

	actionCmd := c.bldr.Build(cmdRendered)

	action = &Action{
		Mode: commandTmpl.Mode,
		cmd:  actionCmd,
	}

	if strings.HasPrefix(cmdRendered, "git clone") {
		action.OnFinished(func() { c.rfs.Refresh(repo) })
	}

	return action, true
}

func (c *Commander) renderCommandTemplate(repo repos.Repository, command string) (string, error) {
	cloneDir, err := c.rfs.FindCloneDirectory(repo)
	if err != nil {
		return "", err
	}

	return quicktmpl.Render(command, quicktmpl.Data{
		"CloneDir": cloneDir,
		"Repo":     repo,
	})
}
