package repos

import (
	"context"
)

type Repository struct {
	ID          int
	RemoteID    string
	Name        string
	Username    string
	Description string
	CloneURL    string
	CloneSSHURL string
	IsFork      bool
}

type RepositoryClient interface {
	GetAllByUsername(ctx context.Context, username string) ([]Repository, error)
	GetOneByUsername(ctx context.Context, username, name string) (Repository, error)

	// GetReadme returns the README.md content of the repostiroy if it's present.
	// If the README.md is not present, it returns an empty string.
	GetReadme(ctx context.Context, username, name string) (string, error)
}
