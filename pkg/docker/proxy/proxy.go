package dockerproxy

import (
	"github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/docker/client"
)

type DockerProxyHandler struct {
	repoStore binaryrepo.RepoStore
	fs 		binaryrepo.Filestore
	registryClient client.RegistryClient
}

func NewDockerProxyHandler(repoStore binaryrepo.RepoStore,
						   fs binaryrepo.Filestore,
						   registryClient client.RegistryClient) * DockerProxyHandler {
	return &DockerProxyHandler {
		repoStore: repoStore,
		fs: fs,
		registryClient: registryClient,
	}
}
