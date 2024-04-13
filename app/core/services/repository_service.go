package services

import (
	"context"
	"database/sql"
	"sync"

	"github.com/hay-kot/repomgr/app/core/bus"
	"github.com/hay-kot/repomgr/app/core/db"
	"github.com/hay-kot/repomgr/app/core/db/migrations"
	"github.com/hay-kot/repomgr/app/repos"
)

type RepositoryService struct {
	sql         *sql.DB
	db          *db.Queries
	mu          sync.RWMutex
	clonedrepos map[string]bool
}

func NewRepositoryService(s *sql.DB, b *bus.EventBus) (*RepositoryService, error) {
	_, err := s.Exec(migrations.Schema)
	if err != nil {
		return nil, err
	}

	rs := &RepositoryService{
		sql:         s,
		db:          db.New(s),
		clonedrepos: make(map[string]bool),
	}

	b.SubCloneEvent(func(rce bus.RepoClonedEvent) {
		rs.SetCloned(rce.Repo, true)
	})

	return rs, nil
}

func (s *RepositoryService) SetCloned(repo repos.Repository, v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clonedrepos[repo.CloneURL] = v
}

func (s *RepositoryService) IsCloned(repo repos.Repository) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.clonedrepos[repo.CloneURL]
	if !ok {
		return false
	}

	return v
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
			HTMLURL:     item.HtmlUrl,
			CloneURL:    item.CloneUrl,
			CloneSSHURL: item.CloneSshUrl,
			IsFork:      item.IsFork,
			ForkURL:     item.ForkUrl,
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

func (s *RepositoryService) UpsertOne(ctx context.Context, item repos.Repository) error {
	return s.UpsertMany(ctx, []repos.Repository{item})
}

func (s *RepositoryService) GetReadme(ctx context.Context, repoID int) ([]byte, error) {
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
