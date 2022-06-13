package dockerregistry

import (
	binaryrepo "github.com/martencassel/binaryrepo"
)

type DockerRegistryHandler struct {
	rs binaryrepo.RepoStore
	fs 		binaryrepo.Filestore
	ts binaryrepo.TagStore
	uploader binaryrepo.Uploader
}

func NewDockerRegistryHandler(rs binaryrepo.RepoStore,
						      fs binaryrepo.Filestore,
							  ts binaryrepo.TagStore,
							  uploader binaryrepo.Uploader) * DockerRegistryHandler {
	return &DockerRegistryHandler {
		rs: rs,
		fs: fs,
		ts: ts,
		uploader: uploader,
	}
}
