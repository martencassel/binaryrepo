package registry

import (
	"net/http"

	log "github.com/rs/zerolog/log"
)

// PathPutManifest URL.
const PathPutManifest = "/repo/{repo-name}/v2/{name}/manifests/{reference}"

// Put the manifest identified by name and reference where reference can be a tag or digest.
func (registry *DockerRegistry) PutManifest(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.putManifest %s %s", r.Method, r.URL.Path)
}

func (registry *DockerRegistry) HeadManifestHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.checkManifest %s %s", r.Method, r.URL.Path)
}

func (registry *DockerRegistry) GetManifestHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.getManifest %s %s", r.Method, r.URL.Path)
}
