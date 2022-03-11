package registry

import (
	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/docker/uploader"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	log "github.com/rs/zerolog/log"
)

const PathChunkedUpload = "/repo/{repo-name}/v2/{name}/blobs/uploads/{uuid}"

func RegisterHandlers(r *mux.Router, fs *filestore.FileStore, repoIndex *repo.RepoIndex, uploader *uploader.UploadManager) {
	log.Info().Msgf("Registering docker registry handlers")
	registry := NewDockerRegistry(fs, repoIndex, uploader)
	r.HandleFunc(RegistryPathVersion, registry.VersionHandler).Methods("GET")
}
