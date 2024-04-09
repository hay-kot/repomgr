// Package commands contains the CLI commands for the application
package commands

import (
	"sync"

	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/repos"
)

type Controller struct {
	Flags *Flags
	cc    clientCache
}

func NewController() *Controller {
	return &Controller{
		Flags: &Flags{},
		cc: clientCache{
			cache: make(map[cacheKey]repos.RepositoryClient),
		},
	}
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
