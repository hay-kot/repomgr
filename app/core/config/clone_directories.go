package config

import (
	"fmt"
	"path/filepath"
)

type CloneDirectories struct {
	Default  string    `toml:"default"`
	Matchers []Matcher `toml:"matchers"`
}

func (c CloneDirectories) FindMatch(str string) string {
	for _, matcher := range c.Matchers {
		if matcher.IsMatch(str) {
			return matcher.Directory
		}
	}

	return c.Default
}

type Matcher struct {
	Match     string `toml:"match"`
	Directory string `toml:"dir"`
}

func (m Matcher) IsMatch(str string) bool {
	ok, _ := filepath.Match(m.Match, str)
	return ok
}

func (c CloneDirectories) Validate() error {
	if c.Default == "" {
		return fmt.Errorf("default clone directory is required")
	}

	for i, matcher := range c.Matchers {
		if matcher.Match == "" {
			return fmt.Errorf("match is required for clone directory matcher %d", i)
		}

		_, err := filepath.Match(matcher.Match, "")
		if err != nil {
			return fmt.Errorf("invalid match pattern for clone directory matcher %d: %w", i, err)
		}

		if matcher.Directory == "" {
			return fmt.Errorf("directory is required for clone directory matcher %d", i)
		}
	}

	return nil
}
