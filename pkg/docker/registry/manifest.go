package registry

import (
	"net/http"

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
}

func (registry *DockerRegistry) HasManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.hasManifestHandler %s %s", req.Method, req.URL.Path)
}

func (registry *DockerRegistry) GetManifestHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.getManifest %s %s", req.Method, req.URL.Path)
}
