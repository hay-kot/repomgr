package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hay-kot/repomgr/app/core/db"
	"github.com/hay-kot/repomgr/app/repos"
)

type ArtifactType string

func (a ArtifactType) String() string {
	return string(a)
}

const (
	ArtifactTypeReadme ArtifactType = "repo.readme"
)

type RepositoryService struct {
	sql *sql.DB
	db  *db.Queries
}

func NewRepositoryService(s *sql.DB) *RepositoryService {
	return &RepositoryService{
		sql: s,
		db:  db.New(s),
	}
}

func (s *RepositoryService) GetAll(ctx context.Context) ([]repos.Repository, error) {
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
			CloneURL:    item.CloneUrl,
			CloneSSHURL: item.CloneSshUrl,
			IsFork:      item.IsFork,
		}
	}

	return results, nil
}

func (s *RepositoryService) UpsertMany(ctx context.Context, items []repos.Repository) error {
	// TODO: implement transactions
	tx := s.db
	for _, item := range items {
		_, err := tx.RepoUpsert(ctx, db.RepoUpsertParams{
			RemoteID:    item.RemoteID,
			Name:        item.Name,
			Username:    item.Username,
			Description: item.Description,
			CloneUrl:    item.CloneURL,
			CloneSshUrl: item.CloneSSHURL,
			IsFork:      item.IsFork,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *RepositoryService) UpsertOne(ctx context.Context, item repos.Repository) error {
	return s.UpsertMany(ctx, []repos.Repository{item})
}

var ErrNoReadmeFound = errors.New("no readme found")

func (s *RepositoryService) GetReadme(ctx context.Context, repoID int) ([]byte, error) {
	v, err := s.db.RepoArtifactByType(ctx, db.RepoArtifactByTypeParams{
		RepositoryID: int64(repoID),
		DataType:     ArtifactTypeReadme.String(),
	})
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, ErrNoReadmeFound
	}

	return v[0].Data, nil
}

func (s *RepositoryService) SetReadme(ctx context.Context, repoID int, data []byte) error {
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
