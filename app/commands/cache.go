package commands

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/repomgr/app/commands/ui"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/repos"
)

func (ctrl *Controller) Cache(ctx context.Context) error {
	msgch := make(chan string)
	defer close(msgch)
	spinner := ui.NewSpinner(msgch, "Caching repositories...")

	sem := make(chan struct{}, ctrl.conf.Concurrency)
	defer close(sem)

	wg := &sync.WaitGroup{}
	wg.Add(len(ctrl.conf.Sources))

	total := 0
	appendTotal := func(v int) {
		total += v
		msgch <- fmt.Sprintf("Total repositories: %d", total)
	}

	for _, source := range ctrl.conf.Sources {
		go func(source config.Source) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			client, err := ctrl.client(source.Type, source.Token())
			if err != nil {
				return
			}

			repos, err := client.GetAllByUsername(ctx, source.Username)
			if err != nil {
				return
			}

			appendTotal(len(repos))
		}(source)
	}

	wg.Wait()
	p := tea.NewProgram(spinner)
	_, err := p.Run()
	if err != nil {
		panic(err)
	}
	return nil
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
