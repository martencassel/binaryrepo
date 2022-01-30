package registry

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/rs/zerolog/log"
)

// PathPutManifest URL.
const RegistryPathVersion = "/repo/{repo-name}/v2/"

// VersionHandler implements GET baseURL/repo/v2/
func (registry *DockerRegistry) VersionHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.version %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	repo := registry.index.FindRepo(repoName)
	if repo == nil {
		log.Printf("Repo %s was not found", repoName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
