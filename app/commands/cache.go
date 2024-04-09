package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/repos"
)

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

func (ctrl *Controller) Cache(ctx context.Context, cfg *config.Config) error {
	for _, source := range cfg.Sources {
		client, err := ctrl.client(source.Type, source.Token())
		if err != nil {
			return err
		}

		repos, err := client.GetAllByUsername(ctx, source.Username)
		if err != nil {
			return err
		}

		// TODO: storage mechanism

		// write to repos.json
		f, err := os.Create("repos.json")
		if err != nil {
			return err
		}

		defer f.Close()

		err = json.NewEncoder(f).Encode(repos)
		if err != nil {
			return err
		}
	}

	return nil
}
