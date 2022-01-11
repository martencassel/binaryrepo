package dockerproxy

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

// RegisterHandlers is a method that registers
// all the docker proxy paths.
func RegisterHandlers(r *mux.Router, fs *filestore.FileStore, repoIndex *repo.RepoIndex) {
	log.Print("Registering docker proxy handlers")
	p := NewDockerProxyApp()
	p.index = repoIndex
	r.HandleFunc(PathVersionUrl1, p.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathVersionUrl2, p.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathHeadManifest1, p.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathHeadManifest2, p.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathGetManifest1, p.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetManifest2, p.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob1, p.GetBlobHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob2, p.GetBlobHandler).Methods(http.MethodGet)
	r.HandleFunc(PathServeBlobURL, p.ServeBlobHandler).Methods(http.MethodGet)
}
