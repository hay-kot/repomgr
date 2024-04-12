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
	username := ""
	if repo.GetOwner() != nil {
		username = repo.GetOwner().GetLogin()
	} else if repo.GetOrganization() != nil {
		username = repo.GetOrganization().GetLogin()
	}

	fork_url := ""
	if repo.GetFork() {
		parenet := repo.GetParent()
		if parenet != nil {
			fork_url = parenet.GetHTMLURL()
		} else {
			log.Warn().
				Str("repo", repo.GetHTMLURL()).
				Msg(" forked repo does not have parent")
		}
	}

	return Repository{
		RemoteID:    strconv.FormatInt(repo.GetID(), 10),
		Name:        repo.GetName(),
		Username:    username,
		Description: repo.GetDescription(),
		HTMLURL:     repo.GetHTMLURL(),
		CloneURL:    repo.GetCloneURL(),
		CloneSSHURL: repo.GetSSHURL(),
		IsFork:      repo.GetFork(),
		ForkURL:     fork_url,
	}
}

// GetAllByUsername implements RepositoryClient.
func (g *GithubClient) GetAllByUsername(ctx context.Context, username string) ([]Repository, error) {
	opt := &github.RepositoryListByAuthenticatedUserOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 200},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := g.client.Repositories.ListByAuthenticatedUser(ctx, opt)
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
