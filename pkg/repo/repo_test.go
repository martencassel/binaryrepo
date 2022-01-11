package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoInit(t *testing.T) {
	index := NewRepoIndex()
	repo1 := Repo{
		ID:      1,
		Name:    "docker-remote-1",
		Type:    Remote,
		PkgType: Docker,
		URL:     "https://registry-1.docker.io",
	}
	index.AddRepo(repo1)
	repo2 := Repo{
		ID:      2,
		Name:    "docker-remote-2",
		Type:    Remote,
		PkgType: Docker,
		URL:     "https://registry-1.docker.io",
	}
	index.AddRepo(repo2)
	assert.Equal(t, 2, len(index.GetRepos()))
	_repo := index.FindRepo("docker-remote-1")
	if _repo == nil {
		t.Errorf("Repo not found")
	}
}
