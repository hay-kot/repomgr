package repostore

type ArtifactType string

func (a ArtifactType) String() string {
	return string(a)
}

const (
	ArtifactTypeReadme ArtifactType = "repo.readme"
)
