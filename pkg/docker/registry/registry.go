package registry

import (
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

type DockerRegistry struct {
	fs    *filestore.FileStore
	index *repo.RepoIndex
}

func NewDockerRegistry(fs *filestore.FileStore, index *repo.RepoIndex) *DockerRegistry {
	registry := &DockerRegistry{}
	registry.fs = fs
	registry.index = index
	return registry
}
