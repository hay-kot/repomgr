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
			WithMaxGoroutines(ctrl.conf.Concurrency).
			WithErrors().
			WithContext(ctx)

		total := 0
		appendTotal := func(v int) {
			total += v
			msgch <- fmt.Sprintf("Total repositories: %d", total)
		}

		collectionch := make(chan []repos.Repository, 1)

		for i := range ctrl.conf.Sources {
			source := ctrl.conf.Sources[i]
			wg.Go(func(ctx context.Context) error {
				client, err := ctrl.client(source.Type, source.Token())
				if err != nil {
					return err
				}

				repos, err := client.GetAllByUsername(ctx, source.Username)
				if err != nil {
					return err
				}

				collectionch <- repos

				appendTotal(len(repos))
				return nil
			})
		}

		items := make([]repos.Repository, 0, total)
		colwg := pool.New()
		colwg.Go(func() {
			for repos := range collectionch {
				items = append(items, repos...)
			}
		})

		err := wg.Wait()
		if err != nil {
			return err
		}

		close(collectionch)
		colwg.Wait()

		msgch <- fmt.Sprintf("Total repositories: %d", len(items))
		msgch <- "Saving repositories to database..."

		return ctrl.repos.UpsertMany(ctx, items)
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
