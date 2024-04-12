package bus

import "github.com/hay-kot/repomgr/app/repos"

type Topic string

var TopicRepoCloned Topic = "repo.cloned"

type eventData struct {
	topic Topic
	data  any
}

type RepoClonedEvent struct {
	Repo     repos.Repository
	CloneDir string
}
