package registry

import (
	tagstore "github.com/martencassel/binaryrepo/pkg/docker/tagstore"
	"github.com/martencassel/binaryrepo/pkg/docker/uploader"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

type DockerRegistry struct {
	fs       *filestore.FileStore
	index    *repo.RepoIndex
	tagstore *tagstore.TagStore
	uploader *uploader.UploadManager
}

func NewDockerRegistry(fs *filestore.FileStore, index *repo.RepoIndex, uploader *uploader.UploadManager) *DockerRegistry {
	registry := &DockerRegistry{}
	registry.fs = fs
	registry.index = index
	registry.tagstore = tagstore.NewTagStore("/tmp/tagstore")
	registry.uploader = uploader
	return registry
}
