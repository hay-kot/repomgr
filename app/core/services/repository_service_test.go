package services

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/hay-kot/repomgr/app/core/db/migrations"
	"github.com/hay-kot/repomgr/app/repos"
	"github.com/matryer/is"
	_ "modernc.org/sqlite"
)

func tServiceFactory(t *testing.T) *RepositoryService {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(migrations.Schema)
	if err != nil {
		t.Fatal(err)
	}

	return NewRepositoryService(db)
}

func factory(n int) []repos.Repository {
	results := make([]repos.Repository, n)
	for n := range results {
		results[n] = repos.Repository{
			RemoteID:    faker.UUIDHyphenated(),
			Name:        faker.Name(),
			Username:    faker.Username(),
			Description: faker.Sentence(),
			CloneURL:    faker.URL(),
			CloneSSHURL: faker.URL(),
			IsFork:      false,
		}
	}

	return results
}

func compareRepository(is *is.I, got, want repos.Repository) {
	is.Helper()
	is.Equal(got.RemoteID, want.RemoteID)
	is.Equal(got.Name, want.Name)
	is.Equal(got.Username, want.Username)
	is.Equal(got.Description, want.Description)
	is.Equal(got.CloneURL, want.CloneURL)
	is.Equal(got.CloneSSHURL, want.CloneSSHURL)
	is.Equal(got.IsFork, want.IsFork)
}

func Test_RepositoryService_UpsertMany(t *testing.T) {
	const Count = 20

	service := tServiceFactory(t)

	tocreate := factory(Count)

	is := is.New(t)
	err := service.UpsertMany(context.Background(), tocreate)
	is.NoErr(err)

	all, err := service.GetAll(context.Background())
	is.NoErr(err)

	is.Equal(len(all), Count) // 20 records should be created

	// re-insert the same records
	err = service.UpsertMany(context.Background(), tocreate)
	is.NoErr(err)

	all, err = service.GetAll(context.Background())
	is.NoErr(err)

	is.Equal(len(all), Count) // 20 records should be created

	// validate records
	for _, got := range all {
		for _, want := range tocreate {
			if got.RemoteID == want.RemoteID {
				compareRepository(is, got, want)
			}
		}
	}
}

func Test_RepositoryService_UpsertOne(t *testing.T) {
	service := tServiceFactory(t)

	is := is.New(t)

	item := factory(1)[0]
	err := service.UpsertOne(context.Background(), item)
	is.NoErr(err)

	all, err := service.GetAll(context.Background())

	is.NoErr(err)
	if len(all) != 1 {
		t.Fatalf("expected 1 record, got %d", len(all))
	}

  compareRepository(is, all[0], item) 
}
