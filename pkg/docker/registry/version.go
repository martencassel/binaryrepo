package dockerregistry

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/rs/zerolog/log"
)

/*
	GET /v2/

	If a 200 OK response, registry implements V2.(1) registry API.
	If a 401 Unauthorized, the client should take action based on the contents of "WWW-Authenticate" header and try this endpoint again.
	If 404 Not Found, client should proceed with the assumption that the registry does not implement V2 of the API.
	When a 200 OK or 401 unauthorized response is returned, "Docker-Distribution-API-Version" headers should be set to "registry/2.0"

*/

// PathVersion URL.
const PathVersionUrl1 = "/repo/{repo-name}/v2"
const PathVersionUrl2 = "/repo/{repo-name}/v2/"

// VersionHandler implements GET baseURL/repo/v2/
func (p *DockerRegistryHandler) VersionHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerRegistry.VersionHandler: %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	log.Info().Msgf("dockerRegistry.VersionHandler: Repository: %s", repoName)
	rw.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
	rw.WriteHeader(http.StatusOK)
}