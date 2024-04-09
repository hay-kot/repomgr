package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hay-kot/repomgr/app/commands/ui"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/sourcegraph/conc/pool"
)

func (ctrl *Controller) Cache(ctx context.Context) error {
	return ui.NewSpinnerFunc("Cacheing Repositories...", func(msgch chan<- string) error {
		wg := pool.New().
			WithErrors().
			WithContext(ctx)

		total := 0
		appendTotal := func(v int) {
			total += v
			msgch <- fmt.Sprintf("Total repositories: %d", total)
		}

		sem := make(chan struct{}, ctrl.conf.Concurrency)

		for i := range ctrl.conf.Sources {
			source := ctrl.conf.Sources[i]
			wg.Go(func(ctx context.Context) error {
				sem <- struct{}{}
				defer func() { <-sem }()

				client, err := ctrl.client(source.Type, source.Token())
				if err != nil {
					return err
				}

				repos, err := client.GetAllByUsername(ctx, source.Username)
				if err != nil {
					return err
				}

				appendTotal(len(repos))
				return nil
			})
		}

		return wg.Wait()
	})
}

func (ctrl *Controller) client(t config.SourceType, token string) (repos.RepositoryClient, error) {
	if client, ok := ctrl.cc.get(t, token); ok {
		return client, nil
	}

	var client repos.RepositoryClient
	switch t {
	case config.SourceTypeGithub:
		client = repos.NewGithubClient(http.DefaultClient, token)
	default:
		return nil, fmt.Errorf("unsupported repository source type: %s", t)
	}

	ctrl.cc.set(t, token, client)
	return client, nil
}
