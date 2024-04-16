package services

import (
	"fmt"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/rs/zerolog/log"
)

type CommandResult struct {
	IsExit      bool
	ExitMessage string
}

func (s *AppService) RunInteractive(repo repos.Repository, command string) (*CommandHandle, error) {
	cmd, args, err := s.prepareCommand(repo, command)
	if err != nil {
		return nil, err
	}

	return s.exec.ExecuteHandler(cmd, args...)
}

func (s *AppService) Run(repo repos.Repository, command string) (CommandResult, error) {
	cmd, args, err := s.prepareCommand(repo, command)
	if err != nil {
		return CommandResult{}, err
	}

	// Special Case
	if strings.HasPrefix(command, string(AppCommandExit)) {
		return CommandResult{
			IsExit:      true,
			ExitMessage: strings.TrimPrefix(cmd, ":Exit"),
		}, nil
	}

	err = s.exec.Execute(cmd, args...)
	if err != nil {
		return CommandResult{}, err
	}

	switch {
	case strings.HasPrefix(command, "git clone"):
		log.Debug().Str("repo", repo.DisplayName()).Msg("emitting clone event")
		cloneDir, err := s.findCloneDirectory(repo)
		if err != nil {
			return CommandResult{}, err
		}

		s.bus.PubCloneEvent(repo, cloneDir)
		s.clonecache.Set(repo.CloneURL, true)
	}

	return CommandResult{}, nil
}

func (s *AppService) GetBoundCommand(cmd tea.KeyType) (bool, string) {
	ok, cmdstr := TranslateTeaKey(cmd)
	if !ok {
		log.Debug().Str(" key", cmdstr).Msg("key not found")
		return false, ""
	}

	for key, command := range s.cfg.KeyBindings {
		if key == cmdstr {
			log.Debug().Str("key", key).Str("command", command.String()).Msg("found key")
			return true, command.String()
		}
	}

	return false, ""
}

func (s *AppService) renderCommandTemplate(repo repos.Repository, command string) (string, error) {
	tmpl, err := template.New("command").Parse(command)
	if err != nil {
		return "", err
	}

	cloneDir, err := s.findCloneDirectory(repo)
	if err != nil {
		return "", err
	}

	b := &strings.Builder{}
	err = tmpl.Execute(b, map[string]any{
		"CloneDir": cloneDir,
		"Repo":     repo,
	})
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (s *AppService) prepareCommand(repo repos.Repository, command string) (string, []string, error) {
	if strings.HasPrefix(command, ":") {
		// special command
		switch {
		case strings.HasPrefix(command, AppCommandExit.String()):
			cmd := strings.TrimPrefix(command, AppCommandExit.String()+" ")

			cmd, err := s.renderCommandTemplate(repo, cmd)
			if err != nil {
				return "", nil, err
			}

			return cmd, nil, nil
		case strings.HasPrefix(command, AppCommandFork.String()):
			// TODO implement fork command
			panic("not implemented")
		default:
			return "", nil, fmt.Errorf("unknown command: '%s'", command)
		}
	}

	command, err := s.renderCommandTemplate(repo, command)
	log.Debug().Str("command", command).Msg("rendered command")
	if err != nil {
		return "", nil, err
	}

	args := splitWithQuotes(command)
	if len(args) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}

	cmd, args := args[0], args[1:]
	return cmd, args, nil
}

func splitWithQuotes(input string) []string {
	var parts []string
	var currentPart strings.Builder
	insideQuotes := false

	for _, char := range input {
		if char == '\'' {
			insideQuotes = !insideQuotes
		} else if char == ' ' && !insideQuotes {
			// If it's a space and we're not inside quotes,
			// we consider it as a separator between parts.
			parts = append(parts, currentPart.String())
			currentPart.Reset()
			continue
		}
		// Append the character to the current part.
		currentPart.WriteRune(char)
	}

	// Append the last part after the loop.
	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}

type AppCommand string

func (c AppCommand) String() string {
	return string(c)
}

const (
	AppCommandFork AppCommand = ":GitFork"
	AppCommandExit AppCommand = ":Exit"
)
