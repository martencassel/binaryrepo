package dockerproxy

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/rs/zerolog/log"
)

// PathVersion URL.
const PathVersionUrl1 = "/repo/{repo-name}/v2"
const PathVersionUrl2 = "/repo/{repo-name}/v2/"

// VersionHandler implements GET baseURL/repo/v2/
func (p *DockerProxyApp) VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("VersionHandler: %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	log.Printf("Repo Name: %s\n", repoName)
	repo := p.index.FindRepo(repoName)
	log.Info().Msgf("%v", r)
	if repo == nil {
		log.Printf("Repo %s was not found", repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
