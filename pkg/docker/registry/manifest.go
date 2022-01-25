package registry

import (
	"net/http"
)

// PathPutManifest URL.
const PathPutManifest = "/repo/{repo-name}/v2/{name}/manifests/{reference}"

// Put the manifest identified by name and reference where reference can be a tag or digest.
func (registry *DockerRegistry) PutManifest(w http.ResponseWriter, r *http.Request) {
	////log.Info().Msgf("HeadManifestHandler %s %s", r.Method, r.URL.Path)
}

func (registry *DockerRegistry) HeadManifestHandler(w http.ResponseWriter, r *http.Request) {
	////log.Info().Msgf("HeadManifestHandler %s %s", r.Method, r.URL.Path)
}

func (registry *DockerRegistry) GetManifestHandler(w http.ResponseWriter, r *http.Request) {
	////log.Info().Msgf("GetManifestHandler %s %s", r.Method, r.URL.Path)
}
