package config

import (
	"fmt"
	"os"
	"strings"
)

type SourceType string

var SourceTypeGithub SourceType = "github"

func (st SourceType) String() string {
	return string(st)
}

func (st SourceType) IsValid() bool {
	switch st {
	case SourceTypeGithub:
		return true
	default:
		return false
	}
}

type Source struct {
	Type     SourceType `toml:"type"`
	Username string     `toml:"username"`
	TokenKey string     `toml:"token"`
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

	if !s.Type.IsValid() {
		return fmt.Errorf("source type is invalid")
	}

	if s.Username == "" {
		return fmt.Errorf("source username is required")
	}
	return nil
}
