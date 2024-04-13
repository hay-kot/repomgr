// Package commands contains the CLI commands for the application
package commands

import (
	"context"
	"database/sql"
	"sync"

	"github.com/hay-kot/repomgr/app/core/bus"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/core/services"
	"github.com/hay-kot/repomgr/app/repos"
)

type Controller struct {
	app  *services.AppService
	conf *config.Config
	cc   clientCache
	bus  *bus.EventBus
}

func NewController(conf *config.Config) (*Controller, error) {
	bus := bus.NewEventBus(10)
	go func() {
		_ = bus.Start(context.Background())
	}()

	sqldb, err := sql.Open("sqlite", conf.Database.DNS())
	if err != nil {
		return nil, err
	}
	exec := services.NewShellExecutor(conf.Shell)

	// defer sqldb.Close()
	app, err := services.NewAppService(sqldb, conf, exec, bus)
	if err != nil {
		return nil, err
	}

	return &Controller{
		bus:  bus,
		conf: conf,
		app:  app,
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
