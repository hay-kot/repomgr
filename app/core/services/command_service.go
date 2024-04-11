package services

import (
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/rs/zerolog/log"
)

type Executor interface {
	Execute(cmd string, args ...string) error
}

type ShellExecutor struct{}

func (e ShellExecutor) Execute(cmd string, args ...string) error {
	log.Debug().Str("cmd", cmd).Strs("args", args).Msg("executing command")
	return exec.Command(cmd, args...).Run()
}

type CommandService struct {
	dirs config.CloneDirectories
	keys config.KeyBindings
	exec Executor
}

func NewCommandService(
	dirs config.CloneDirectories,
	keys config.KeyBindings,
	e Executor,
) *CommandService {
	return &CommandService{
		dirs: dirs,
		keys: keys,
		exec: e,
	}
}

type AppCommand string

const (
	AppCommandClone AppCommand = ":GitClone"
	AppCommandFork  AppCommand = ":GitFork"
	AppCommandPull  AppCommand = ":GitPull"
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
		case AppCommandClone:
			command = "git clone {{ .Repo.CloneSSHURL }} {{ .CloneDir }}"
		case AppCommandFork:
			// TODO implement fork command
			panic("not implemented")
		case AppCommandPull:
			command = "git pull {{ .CloneDir }}"
		default:
			return "", nil, fmt.Errorf("unknown command: '%s'", command)
		}
	}

	command, err := s.renderCommandTemplate(repo, command)
	if err != nil {
		return "", nil, err
	}

	args := strings.Split(command, " ")
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

	return s.exec.Execute(cmd, args...)
}
