package services

import (
	"database/sql"
	"os"
	"strings"
	"text/template"

	"github.com/hay-kot/repomgr/app/core/bus"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/core/db"
	"github.com/hay-kot/repomgr/app/core/db/migrations"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/hay-kot/repomgr/internal/cache"
	"github.com/rs/zerolog/log"
)

type AppService struct {
	exec       Executor
	cfg        *config.Config
	sql        *sql.DB
	db         *db.Queries
	bus        *bus.EventBus
	clonecache cache.Cache[bool]
}

func NewAppService(s *sql.DB, cfg *config.Config, exec Executor, b *bus.EventBus) (*AppService, error) {
	_, err := s.Exec(migrations.Schema)
	if err != nil {
		return nil, err
	}

	app := &AppService{
		cfg:  cfg,
		sql:  s,
		exec: exec,
		db:   db.New(s),
		bus:  b,
	}

	return app, nil
}

func (s *AppService) findCloneDirectory(repo repos.Repository) (string, error) {
	dirtmpl := s.cfg.CloneDirectories.FindMatch(repo.DisplayName())

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

func (s *AppService) IsCloned(r repos.Repository) bool {
	if v, ok := s.clonecache.Get(r.CloneURL); ok {
		return v
	}

	path, err := s.findCloneDirectory(r)
	if err != nil {
		log.Warn().Err(err).Msg("failed to find clone directory")
		return false
	}

	// check if directory exists
	if _, err := os.Stat(path); err == nil {
		s.clonecache.Set(r.CloneURL, true)
	} else {
		s.clonecache.Set(r.CloneURL, false)
	}

	return false
}
