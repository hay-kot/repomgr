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

func (s *RepositoryService) GetReadme(ctx context.Context, repoID int) ([]byte, error) {
	v, err := s.db.RepoArtifactByType(ctx, db.RepoArtifactByTypeParams{
		RepositoryID: int64(repoID),
		Type:         ArtifactTypeReadme.String(),
	})
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, errors.New("no readme found")
	}

	return v[0].Data, nil
}

func (s *RepositoryService) SetReadme(ctx context.Context, repoID int, data []byte) error {
	_, err := s.db.RepoCreateArtifact(ctx, db.RepoCreateArtifactParams{
		RepositoryID: int64(repoID),
		Type:         ArtifactTypeReadme.String(),
		Data:         data,
	})
	if err != nil {
		return err
	}

	return err
}
