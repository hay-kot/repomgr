package config

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
)

type SourceType string

var SourceTypeGithub SourceType = "github"

type Source struct {
	Type          SourceType `toml:"type"`
	Username      string     `toml:"username"`
	Organizations []string   `toml:"organizations"`
	TokenKey      string     `toml:"token"`
}

func (s Source) Token() string {
	if strings.HasPrefix(s.TokenKey, "env:") {
		return os.Getenv(strings.TrimPrefix(s.TokenKey, "env:"))
	}

	return s.TokenKey
}

func (s Source) Validate() error {
	if s.Type == "" {
		return fmt.Errorf("source type is required")
	}

	types := []SourceType{SourceTypeGithub}
	var found bool
	for _, t := range types {
		if t == s.Type {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("source type is invalid")
	}

	if s.Username == "" {
		return fmt.Errorf("source username is required")
	}
	return nil
}

type Config struct {
	ProjectDir  string      `toml:"project_dir"`
	KeyBindings KeyBindings `toml:"key_bindings"`
	Concurrency int         `toml:"concurrency"`
	Sources     []Source    `toml:"sources"`
}

func New(reader io.Reader) (*Config, error) {
	cfg := Config{
		Concurrency: runtime.NumCPU(),
	}

	_, err := toml.NewDecoder(reader).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c Config) Validate() error {
	if c.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}

	if len(c.Sources) == 0 {
		return fmt.Errorf("sources are required")
	}

	if c.ProjectDir == "" {
		return fmt.Errorf("project_dir is required")
	}

	for _, source := range c.Sources {
		if err := source.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c Config) Dump() (string, error) {
	var b strings.Builder
	enc := toml.NewEncoder(&b)
	err := enc.Encode(c)
	return b.String(), err
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
