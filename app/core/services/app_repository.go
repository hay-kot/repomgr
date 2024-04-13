package services

import (
	"context"
	"errors"

	"github.com/hay-kot/repomgr/app/core/db"
	"github.com/hay-kot/repomgr/app/repos"
)

var ErrNoReadmeFound = errors.New("no readme found")

type ArtifactType string

func (a ArtifactType) String() string {
	return string(a)
}

const (
	ArtifactTypeReadme ArtifactType = "repo.readme"
)

func (s *AppService) GetAll(ctx context.Context) ([]repos.Repository, error) {
	v, err := s.db.ReposGetAll(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]repos.Repository, len(v))
	for i, item := range v {
		results[i] = repos.Repository{
			ID:          int(item.ID),
			RemoteID:    item.RemoteID,
			Name:        item.Name,
			Username:    item.Username,
			Description: item.Description,
			HTMLURL:     item.HtmlUrl,
			CloneURL:    item.CloneUrl,
			CloneSSHURL: item.CloneSshUrl,
			IsFork:      item.IsFork,
			ForkURL:     item.ForkUrl,
		}
	}

	return results, nil
}

func (s *AppService) UpsertMany(ctx context.Context, items []repos.Repository) error {
	// TODO: implement transactions
	tx := s.db
	for _, item := range items {
		_, err := tx.RepoUpsert(ctx, db.RepoUpsertParams{
			RemoteID:    item.RemoteID,
			Name:        item.Name,
			Username:    item.Username,
			Description: item.Description,
			HtmlUrl:     item.HTMLURL,
			CloneUrl:    item.CloneURL,
			CloneSshUrl: item.CloneSSHURL,
			IsFork:      item.IsFork,
			ForkUrl:     item.ForkURL,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppService) UpsertOne(ctx context.Context, item repos.Repository) error {
	return s.UpsertMany(ctx, []repos.Repository{item})
}

func (s *AppService) GetReadme(ctx context.Context, repoID int) ([]byte, error) {
	v, err := s.db.RepoArtifactByType(ctx, db.RepoArtifactByTypeParams{
		RepositoryID: int64(repoID),
		DataType:     ArtifactTypeReadme.String(),
	})
	if err != nil {
		return nil, err
	}

	if len(v) == 0 || len(v[0].Data) == 0 {
		return nil, ErrNoReadmeFound
	}

	return v[0].Data, nil
}

func (s *AppService) SetReadme(ctx context.Context, repoID int, data []byte) error {
	_, err := s.db.RepoUpsertArtifact(ctx, db.RepoUpsertArtifactParams{
		RepositoryID: int64(repoID),
		DataType:     ArtifactTypeReadme.String(),
		Data:         data,
	})
	if err != nil {
		return err
	}

	return err
}
