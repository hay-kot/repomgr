package services

import (
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/core/bus"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/rs/zerolog/log"
)

type Executor interface {
	Execute(cmd string, args ...string) error
}

type ShellExecutor struct {
	shell string
}

func NewShellExecutor(shell string) ShellExecutor {
	return ShellExecutor{shell: shell}
}

func (e ShellExecutor) Execute(cmd string, args ...string) error {
	log.Debug().Str("cmd", cmd).Strs("args", args).Msg("executing command")
	err := exec.Command("bash", "-c", cmd+" "+strings.Join(args, " ")).Run()
	log.Debug().Err(err).Msg("command executed")
	return err
}

type CommandService struct {
	dirs config.CloneDirectories
	keys config.KeyBindings
	exec Executor
	bus  *bus.EventBus
}

func NewCommandService(
	dirs config.CloneDirectories,
	keys config.KeyBindings,
	e Executor,
	b *bus.EventBus,
) *CommandService {
	return &CommandService{
		dirs: dirs,
		keys: keys,
		exec: e,
		bus:  b,
	}
}

type AppCommand string

const (
	AppCommandFork AppCommand = ":GitFork"
)

func (s *CommandService) GetBoundCommand(cmd tea.KeyType) (bool, string) {
	ok, cmdstr := TranslateTeaKey(cmd)
	if !ok {
		log.Debug().Str(" key", cmdstr).Msg("key not found")
		return false, ""
	}

	for i := range s.keys {
		log.Debug().Str("key", i).Msg("key found")
	}

	for key, command := range s.keys {
		log.Debug().Str("key", key).Str("command", command.String()).Msg("checking key")
		if key == cmdstr {
			log.Debug().Str("key", key).Str("command", command.String()).Msg("found key")
			return true, command.String()
		}
	}

	return false, ""
}

func (s *CommandService) findCloneDirectory(repo repos.Repository) (string, error) {
	dirtmpl := s.dirs.FindMatch(repo.DisplayName())

	tmpl, err := template.New("dir").Parse(dirtmpl)
	if err != nil {
		return "", err
	}

	b := &strings.Builder{}
	err = tmpl.Execute(b, map[string]any{"Repo": repo})
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (s *CommandService) renderCommandTemplate(repo repos.Repository, command string) (string, error) {
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

func (s *CommandService) prepareCommand(repo repos.Repository, command string) (string, []string, error) {
	if strings.HasPrefix(command, ":") {
		// special command
		switch AppCommand(command) {
		case AppCommandFork:
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

func (s *CommandService) Run(repo repos.Repository, command string) error {
	cmd, args, err := s.prepareCommand(repo, command)
	if err != nil {
		return err
	}

	err = s.exec.Execute(cmd, args...)
	if err != nil {
		return err
	}

	if strings.HasPrefix(command, "git clone") {
		log.Debug().Str("repo", repo.DisplayName()).Msg("emitting clone event")
		cloneDir, err := s.findCloneDirectory(repo)
		if err != nil {
			return err
		}

		s.bus.PubCloneEvent(repo, cloneDir)
	}

	return nil
}
