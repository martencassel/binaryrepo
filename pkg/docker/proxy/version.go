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
func (p *DockerProxyHandler) VersionHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerproxyapp.VersionHandler: %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	log.Info().Msgf("dockerproxyapp.VersionHandler: repoName: %s", repoName)
	rw.WriteHeader(http.StatusOK)
}