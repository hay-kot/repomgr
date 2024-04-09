package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/hay-kot/repomgr/app/repos"
	_ "modernc.org/sqlite"
)

func tServiceFactory(t *testing.T) *RepositoryService {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	return NewRepositoryService(db)
}

func Test_RepositoryService_UpsertMany(t *testing.T) {
	service := tServiceFactory(t)

	tocreate := []repos.Repository{
		{
			RemoteID:    "1",
			Name:        "repo1",
			Username:    "mgr-test",
			Description: "test repo",
			CloneURL:    "clone-url",
			CloneSSHURL: "clone-ssh-url",
			IsFork:      false,
		},
	}

	err := service.UpsertMany(context.Background(), tocreate)
	if err != nil {
		t.Errorf("UpsertMany failed: %v", err)
	}
}
