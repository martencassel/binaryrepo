package registry

import (
	tagstore "github.com/martencassel/binaryrepo/pkg/docker/tagstore"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

type DockerRegistry struct {
	fs       *filestore.FileStore
	index    *repo.RepoIndex
	tagstore *tagstore.TagStore
}

func NewDockerRegistry(fs *filestore.FileStore, index *repo.RepoIndex) *DockerRegistry {
	registry := &DockerRegistry{}
	registry.fs = fs
	registry.index = index
	registry.tagstore = tagstore.NewTagStore("/tmp/tagstore")
	return registry
}
