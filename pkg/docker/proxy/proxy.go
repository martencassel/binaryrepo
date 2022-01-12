package dockerproxy

import (
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
)

type DockerProxyApp struct {
	fs    *filestore.FileStore
	index *repo.RepoIndex
}

func NewDockerProxyApp() *DockerProxyApp {
	p := DockerProxyApp{
		fs:    filestore.NewFileStore("/tmp/filestore"),
		index: repo.NewRepoIndex(),
	}
	return &p
}
