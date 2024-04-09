package repos

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/go-github/v61/github"
	"github.com/rs/zerolog/log"
)

var _ RepositoryClient = &GithubClient{}

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(httpclient *http.Client, token string) *GithubClient {
	client := github.NewClient(httpclient).WithAuthToken(token)

	return &GithubClient{client: client}
}

func (g *GithubClient) mapRepository(repo *github.Repository) Repository {
	return Repository{
		RemoteID:    strconv.FormatInt(repo.GetID(), 10),
		Name:        repo.GetName(),
		Username:    repo.GetOwner().GetLogin(),
		Description: repo.GetDescription(),
		CloneURL:    repo.GetCloneURL(),
		CloneSSHURL: repo.GetSSHURL(),
		IsFork:      repo.GetFork(),
	}
}

// GetAllByUsername implements RepositoryClient.
func (g *GithubClient) GetAllByUsername(ctx context.Context, username string) ([]Repository, error) {
	opt := &github.RepositoryListByUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := g.client.Repositories.ListByUser(ctx, username, opt)
		if err != nil {
			log.Err(err).Ctx(ctx).
				Str("username", username).
				Msg("failed to list repositories")
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	results := make([]Repository, len(allRepos))
	for i, repo := range allRepos {
		results[i] = g.mapRepository(repo)
	}

	log.Debug().Ctx(ctx).Int("count", len(results)).Msg("found repositories")
	return results, nil
}

// GetOneByUsername implements RepositoryClient.
func (g *GithubClient) GetOneByUsername(ctx context.Context, username string, name string) (Repository, error) {
	repo, _, err := g.client.Repositories.Get(ctx, username, name)
	if err != nil {
		log.Err(err).Ctx(ctx).
			Str("username", username).
			Str("name", name).
			Msg("failed to get repository")
		return Repository{}, err
	}

	return g.mapRepository(repo), nil
}

func (g *GithubClient) GetReadme(ctx context.Context, username string, name string) (string, error) {
	content, _, err := g.client.Repositories.GetReadme(ctx, username, name, nil)
	if err != nil {
		log.Err(err).Ctx(ctx).
			Str("username", username).
			Str("name", name).
			Msg("failed to get readme")
		return "", err
	}

	return content.GetContent()
}
