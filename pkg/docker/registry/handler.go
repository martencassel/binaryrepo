package registry

import (
	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	log "github.com/rs/zerolog/log"
)

func RegisterHandlers(r *mux.Router, fs *filestore.FileStore, repoIndex *repo.RepoIndex) {
	log.Info().Msgf("Registering docker registry handlers")
	registry := NewDockerRegistry(fs, repoIndex)
	r.HandleFunc(RegistryPathVersion, registry.VersionHandler).Methods("GET")
}
