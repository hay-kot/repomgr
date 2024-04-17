package repofs

import (
	"html/template"
	"os"
	"strings"

	"github.com/hay-kot/repomgr/app/repos"
	"github.com/hay-kot/repomgr/internal/cache"
	"github.com/rs/zerolog/log"
)

type RepoFS struct {
	clonedirs  CloneDirectories
	clonecache cache.Cache[bool]
}

func New(dirs CloneDirectories) *RepoFS {
	return &RepoFS{
		clonedirs:  dirs,
		clonecache: cache.NewMapCache[bool](20),
	}
}

// FindCloneDirectory finds the clone directory for a repository based on the
// CloneDirectories configuration.
func (rfs *RepoFS) FindCloneDirectory(repo repos.Repository) (string, error) {
	dirtmpl := rfs.clonedirs.FindMatch(repo.DisplayName())

	tmpl, err := template.New("dir").Parse(dirtmpl)
	if err != nil {
		return "", err
	}

	b := &strings.Builder{}
	err = tmpl.Execute(b, map[string]any{"Repo": repo})
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

// IsCloned checks if a repository is cloned. These results are cached
// in-memory.
func (rfs *RepoFS) IsCloned(r repos.Repository) bool {
	if v, ok := rfs.clonecache.Get(r.CloneURL); ok {
		return v
	}

	path, err := rfs.FindCloneDirectory(r)
	if err != nil {
		log.Warn().Err(err).Msg("failed to find clone directory")
		return false
	}

	// TODO: replace if virtual file system
	// check if directory exists
	if _, err := os.Stat(path); err == nil {
		rfs.clonecache.Set(r.CloneURL, true)
	} else {
		rfs.clonecache.Set(r.CloneURL, false)
	}

	return false
}

func (rfs *RepoFS) Refresh(repo repos.Repository) error {
	path, err := rfs.FindCloneDirectory(repo)
	if err != nil {
		return err
	}

	// TODO: replace if virtual file system
	// check if directory exists
	if _, err := os.Stat(path); err == nil {
		rfs.clonecache.Set(repo.CloneURL, true)
	} else {
		rfs.clonecache.Set(repo.CloneURL, false)
	}

	return nil
}
