package registry

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/rs/zerolog/log"
)

// PathPutManifest URL.
const PathPutManifest = "/repo/{repo-name}/v2/{name}/manifests/{reference}"

func (registry *DockerRegistry) HasLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.hasLayer %s %s", req.Method, req.URL.Path)
}

// Put the manifest identified by name and reference where reference can be a tag or digest.
func (registry *DockerRegistry) PutManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.putManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]
	reference := vars["reference"]
	b, err := io.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	digest, err := registry.fs.WriteFile(b)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	exists := registry.fs.Exists(digest)
	if !exists {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	// Write manifest file to storage
	digest, err = registry.fs.WriteFile(b)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Docker-Content-Digest", digest.String())
	rw.Header().Set("Location", "/v2/"+name+"/manifests/"+reference)
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	rw.WriteHeader(http.StatusCreated)
}

func (registry *DockerRegistry) HasManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.hasManifestHandler %s %s", req.Method, req.URL.Path)
}

func (registry *DockerRegistry) GetManifestHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.getManifest %s %s", req.Method, req.URL.Path)
}
