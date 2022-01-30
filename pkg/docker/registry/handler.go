package registry

import (
	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

const PathChunkedUpload = "/repo/{repo-name}/v2/{name}/blobs/uploads/{uuid}"

func RegisterHandlers(r *mux.Router, fs *filestore.FileStore, repoIndex *repo.RepoIndex) {
	////log.Info().Msgf("Registering docker registry handlers")
	registry := NewDockerRegistry(fs, repoIndex)
	r.HandleFunc(RegistryPathVersion, registry.VersionHandler).Methods("GET")
}
