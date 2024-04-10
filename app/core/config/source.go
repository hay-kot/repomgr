package config

import (
	"fmt"
	"os"
	"strings"
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
