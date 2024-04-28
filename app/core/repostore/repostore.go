package repostore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hay-kot/repomgr/app/core/db"
	"github.com/hay-kot/repomgr/app/core/db/migrations"
	"github.com/hay-kot/repomgr/app/repos"
)

var ErrNoReadmeFound = errors.New("no readme found")

type RepoStore struct {
	sql *sql.DB
	db  *db.Queries
}

func New(s *sql.DB) (*RepoStore, error) {
	_, err := s.Exec(migrations.Schema)
	if err != nil {
		return nil, err
	}

	return &RepoStore{sql: s, db: db.New(s)}, nil
}

func (s *RepoStore) GetAll(ctx context.Context) ([]repos.Repository, error) {
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
			Owner:       item.Username,
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

func (s *RepoStore) UpsertMany(ctx context.Context, items []repos.Repository) error {
	// TODO: implement transactions
	tx := s.db
	for _, item := range items {
		_, err := tx.RepoUpsert(ctx, db.RepoUpsertParams{
			RemoteID:    item.RemoteID,
			Name:        item.Name,
			Username:    item.Owner,
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

func (s *RepoStore) UpsertOne(ctx context.Context, item repos.Repository) error {
	return s.UpsertMany(ctx, []repos.Repository{item})
}

func (s *RepoStore) GetReadme(ctx context.Context, repoID int) ([]byte, error) {
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

func (s *RepoStore) SetReadme(ctx context.Context, repoID int, data []byte) error {
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
