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

/*
	Existing Layers

	HEAD /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistry) HasLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.HasLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["namespace"]
	repoName := vars["repo-name"]
	d := vars["digest"]
	if req.Method != http.MethodHead {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(d)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !registry.fs.Exists(digest) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := registry.fs.ReadFile(digest)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Docker-Content-Digest", digest.Hex())
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Header().Set("connection", "close")
	rw.WriteHeader(http.StatusOK)
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

/*
	Pulling a layer

	GET /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistry) DownloadLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.DownloadLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]
	repoName := vars["repo-name"]
	d := vars["digest"]
/*	if req.Method != http.MethodHead || req.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}*/
	if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(d)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !registry.fs.Exists(digest) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := registry.fs.ReadFile(digest)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(b)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Docker-Content-Digest", digest.Hex())
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Header().Set("connection", "close")
	rw.WriteHeader(http.StatusOK)
}

/*
	Deleting a Layer

	DELETE /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistry) DeleteLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.DeleteLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]
	repoName := vars["repo-name"]
	d := vars["digest"]
	if req.Method != http.MethodDelete {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(d)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !registry.fs.Exists(digest) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	err = registry.fs.Remove(digest)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
	rw.Header().Set("Content-Length", "0")
	rw.Header().Set("Docker-Content-Digest", digest.Hex())
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
}
