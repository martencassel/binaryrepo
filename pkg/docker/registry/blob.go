package registry

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	digest "github.com/opencontainers/go-digest"
	log "github.com/rs/zerolog/log"
)

// PathInitBlobUpload URL.
const PathInitBlobUpload = "/repo/{repo-name}/v2/{name}/blobs/upload"

// PathExistsBlob URL.
const PathExistsBlob = "/repo/{repo-name}/v2/{name}/blobs/{uuid}"

func (registry *DockerRegistry) ExistsBlob(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.ExistsBlob %s %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	d := vars["digest"]
	if r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !registry.fs.Exists(digest) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := registry.fs.ReadFile(digest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("docker-distribution-api-version", "registry/2.0")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("connection", "close")
	w.WriteHeader(http.StatusAccepted)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (registry *DockerRegistry) DownloadLayer(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.DownloadLayer %s %s", r.Method, r.URL.Path)
	log.Info().Msg("Not implemented")
}

/*
	Deleting a Layer
	DELETE /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistry) DeleteLayer(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.DeleteLayer %s %s", r.Method, r.URL.Path)
}
