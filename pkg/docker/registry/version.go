package registry

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/rs/zerolog/log"
)

// PathPutManifest URL.
const RegistryPathVersion = "/repo/{repo-name}/v2/"

// VersionHandler implements GET baseURL/repo/v2/
func (registry *DockerRegistry) VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.version %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	repo := registry.index.FindRepo(repoName)
	if repo == nil {
		log.Printf("Repo %s was not found", repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)

}
