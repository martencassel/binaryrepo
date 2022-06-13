package dockerregistry

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opencontainers/go-digest"
	"github.com/rs/zerolog/log"
)

/*
	Existing Layers
	HEAD /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistryHandler) ExistingLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.ExistingLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	// name := vars["namespace"]
	// repoName := vars["repo-name"]
	d := vars["digest"]
	log.Info().Msgf("Digest: %s", d)
	if req.Method != http.MethodHead {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
	// 	rw.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	digest, err := digest.Parse(d)
	log.Info().Msg(digest.String())
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

/*
	Get Layers
	GET /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistryHandler) GetLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.GetLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["namespace"]
	repoName := vars["repo-name"]
	d := vars["digest"]
	log.Info().Msgf("Digest: %s", d)
	log.Info().Msgf("Name: %s, Repo: %s\n", name, repoName)
	// if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
	// 	rw.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	digest, err := digest.Parse(d)
	log.Info().Msg(digest.String())
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
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(b))
}

/*
	Delete layer
	DELETE /v2/<name>/blobs/<digest>
*/
func (registry *DockerRegistryHandler) DeleteLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.DeleteLayer %s %s", req.Method, req.URL.Path)

}