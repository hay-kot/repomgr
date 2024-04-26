package repofs

import (
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/hay-kot/repomgr/app/repos"
	"github.com/hay-kot/repomgr/internal/cache"
	"github.com/rs/zerolog/log"
)

type embedfs struct {
	root string
	fs   fs.FS
}

func (e embedfs) Exists(path string) bool {
	// trim root path from path
	path = strings.TrimPrefix(path, e.root)
	_, err := e.fs.Open(path)
	return err == nil
}

type RepoFS struct {
	clonedirs  CloneDirectories
	clonecache cache.Cache[bool]
	fsmap      map[string]embedfs
}

func New(dirs CloneDirectories) *RepoFS {
	fsmap := make(map[string]embedfs, len(dirs.Matchers)+1)

	for _, matcher := range dirs.Matchers {
		dir := matcher.Directory

		// dir should contain template syntax, e.g. "{{.Repo.DisplayName}}"
		// we need to create an fs.FS for each top level directory where the
		// repositories are cloned into
		// if dir = "/path/to/repos/{{.Repo.DisplayName}}" then we need to create
		// an fs.FS for "/path/to/repos"
		prefix := strings.Split(dir, "{{")[0]
		fsmap[dir] = embedfs{
			root: prefix,
			fs:   os.DirFS(prefix),
		}
	}

	// default directory
	defPrefix := strings.Split(dirs.Default, "{{")[0]
	fsmap[dirs.Default] = embedfs{
		root: defPrefix,
		fs:   os.DirFS(defPrefix),
	}

	return &RepoFS{
		clonedirs:  dirs,
		clonecache: cache.NewMapCache[bool](20),
		fsmap:      fsmap,
	}
}

func (rfs *RepoFS) findCloneDirectory(repo repos.Repository) (path string, dirtmpl string, err error) {
	dirtmpl = rfs.clonedirs.FindMatch(repo.DisplayName())

	_, ok := rfs.fsmap[dirtmpl]
	if !ok {
		log.Error().Str("dir", dirtmpl).Msg("no filesystem found for directory")
	}

	tmpl, err := template.New("dir").Parse(dirtmpl)
	if err != nil {
		return "", "", err
	}

	b := &strings.Builder{}
	err = tmpl.Execute(b, map[string]any{"Repo": repo})
	if err != nil {
		return "", "", err
	}

	return b.String(), dirtmpl, nil
}

// FindCloneDirectory finds the clone directory for a repository based on the
// CloneDirectories configuration.
func (rfs *RepoFS) FindCloneDirectory(repo repos.Repository) (string, error) {
	path, _, err := rfs.findCloneDirectory(repo)
	return path, err
}

// IsCloned checks if a repository is cloned. These results are cached
// in-memory.
func (rfs *RepoFS) IsCloned(r repos.Repository) bool {
	if v, ok := rfs.clonecache.Get(r.CloneURL); ok {
		return v
	}

	return rfs.Refresh(r)
}

func (rfs *RepoFS) Refresh(r repos.Repository) (exists bool) {
	path, dirtmpl, err := rfs.findCloneDirectory(r)
	if err != nil {
		log.Warn().Err(err).Msg("failed to find clone directory")
	}

	fs, ok := rfs.fsmap[dirtmpl]
	if !ok {
		log.Error().Str("dir", dirtmpl).Msg("no filesystem found for directory")
	}

	exists = fs.Exists(path)
	rfs.clonecache.Set(r.CloneURL, exists)
	return exists
}
