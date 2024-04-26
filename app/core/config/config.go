package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/hay-kot/repomgr/app/core/commander"
	"github.com/hay-kot/repomgr/app/core/repofs"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	Concurrency      int                     `toml:"concurrency"`
	Shell            string                  `toml:"shell"`
	DotEnvs          []string                `toml:"dotenvs"`
	KeyBindings      commander.KeyBindings   `toml:"key_bindings"`
	Sources          []Source                `toml:"sources"`
	Database         Database                `toml:"database"`
	Logs             Logs                    `toml:"logs"`
	CloneDirectories repofs.CloneDirectories `toml:"clone_directories"`
}

func Default() *Config {
	return &Config{
		Concurrency: runtime.NumCPU(),
		Logs: Logs{
			Level:  zerolog.InfoLevel,
			File:   "",
			Color:  false,
			Format: "text",
		},
		Database: Database{
			File:   "~/config/repomgr/repos.db",
			Params: "_pragma=busy_timeout=2000&_pragma=journal_mode=WAL&_fk=1",
		},
		KeyBindings: commander.NewDefaultKeyBindings(),
	}
}

func New(confpath string, reader io.Reader) (*Config, error) {
	cfg := Default()
	_, err := toml.NewDecoder(reader).Decode(cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	// load dotenvs
	expandedPaths := []string{}
	for _, path := range cfg.DotEnvs {
		expandedPaths = append(expandedPaths, ExpandPath(confpath, path))
	}

	err = godotenv.Load(expandedPaths...)
	if err != nil {
		return nil, err
	}

	cfg.Database.File = ExpandPath(confpath, cfg.Database.File)
	cfg.Logs.File = ExpandPath(confpath, cfg.Logs.File)

	cfg.CloneDirectories.Default = ExpandPath(confpath, cfg.CloneDirectories.Default)
	for i := range cfg.CloneDirectories.Matchers {
		cfg.CloneDirectories.
			Matchers[i].
			Directory = ExpandPath(confpath, cfg.CloneDirectories.Matchers[i].Directory)
	}

	return cfg, nil
}

func (c Config) PrepareDirectories() error {
	dirs := []string{
		filepath.Dir(c.Database.File),
		filepath.Dir(c.Logs.File),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (c Config) Validate() error {
	if c.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0")
	}

	if len(c.Sources) == 0 {
		return fmt.Errorf("sources are required")
	}

	validators := []validator{
		c.KeyBindings,
		c.Database,
		c.CloneDirectories,
	}

	for _, source := range c.Sources {
		validators = append(validators, source)
	}

	for _, v := range validators {
		if err := v.Validate(); err != nil {
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

type Logs struct {
	Level  zerolog.Level `toml:"level"`
	File   string        `toml:"file"`
	Color  bool          `toml:"color"`
	Format string        `toml:"format"`
}

type Database struct {
	File   string `toml:"file"`
	Params string `toml:"params"`
}

func (d Database) Validate() error {
	if d.File == "" {
		return fmt.Errorf("database file is required")
	}

	return nil
}

func (d Database) DNS() string {
	return fmt.Sprintf("file:%s?%s", d.File, d.Params)
}
