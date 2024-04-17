// Package commands contains the CLI commands for the application
package commands

import (
	"database/sql"
	"sync"

	"github.com/hay-kot/repomgr/app/core/commander"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/core/repofs"
	"github.com/hay-kot/repomgr/app/core/repostore"
	"github.com/hay-kot/repomgr/app/repos"
)

type Controller struct {
	rfs       *repofs.RepoFS
	commander *commander.Commander
	store     *repostore.RepoStore
	conf      *config.Config
	cc        clientCache
}

func NewController(conf *config.Config, sqldb *sql.DB) (*Controller, error) {
	rfs := repofs.New(conf.CloneDirectories)

	store, err := repostore.New(sqldb)
	if err != nil {
		return nil, err
	}

	commander := commander.New(conf.KeyBindings, rfs, &commander.ShellCommandBuilder{
		Shell: conf.Shell,
	})

	return &Controller{
		conf:      conf,
		store:     store,
		rfs:       rfs,
		commander: commander,
		cc: clientCache{
			cache: make(map[cacheKey]repos.RepositoryClient),
		},
	}, nil
}

type cacheKey struct {
	clientType  config.SourceType
	clientToken string
}

type clientCache struct {
	mu    sync.RWMutex
	cache map[cacheKey]repos.RepositoryClient
}

func (cc *clientCache) get(t config.SourceType, token string) (repos.RepositoryClient, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	client, ok := cc.cache[cacheKey{clientType: t, clientToken: token}]
	return client, ok
}

func (cc *clientCache) set(t config.SourceType, token string, client repos.RepositoryClient) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	cc.cache[cacheKey{clientType: t, clientToken: token}] = client
}
