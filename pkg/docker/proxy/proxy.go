package dockerproxy

import (
	"context"
	"time"

	regclient "github.com/martencassel/binaryrepo/pkg/docker/client"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
)

type DockerProxyApp struct {
	fs    *filestore.FileStore
	index *repo.RepoIndex
}

func (*DockerProxyApp) NewRegistryClient(domain, username, password, scope, srvaddr string) (*regclient.Registry, error) {
	ctx := context.Background()
	r, err := regclient.New(ctx, regclient.AuthConfig{
		Username:      username,
		Password:      password,
		Scope:         scope,
		ServerAddress: srvaddr,
	}, regclient.Opt{
		Domain:   domain,
		SkipPing: false,
		Timeout:  time.Minute * 10,
		NonSSL:   false,
		Insecure: false,
	})
	return r, err
}

func NewDockerProxyApp() *DockerProxyApp {
	p := DockerProxyApp{
		fs:    filestore.NewFileStore("/tmp/filestore"),
		index: repo.NewRepoIndex(),
	}
	return &p
}
