package repos

import (
	"context"
)

type RepositoryClient interface {
	GetAllByUsername(ctx context.Context, username string) ([]Repository, error)
	GetOneByUsername(ctx context.Context, username, name string) (Repository, error)

	// GetReadme returns the README.md content of the repostiroy if it's present.
	// If the README.md is not present, it returns an empty string.
	GetReadme(ctx context.Context, username, name string) (string, error)
}

type Repository struct {
	ID          int
	RemoteID    string
	Name        string
	Owner       string
	Description string
	HTMLURL     string
	CloneURL    string
	CloneSSHURL string
	IsFork      bool
	ForkURL     string
}

// DisplayName returns the owner and the name of the repository in the format of "owner/name".
func (r Repository) DisplayName() string {
	return r.Owner + "/" + r.Name
}
