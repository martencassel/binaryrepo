package dockerregistry

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type RepoList struct {
	Repositories []string
}

/*
	GET /v2/_catalog
*/
func (registry *DockerRegistryHandler) ListCatalog(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.ListCatalog %s %s", req.Method, req.URL.Path)
//	vars := mux.Vars(req)
//	repos, err := registry.ts.GetRepos()
	// repoList := RepoList{
	// 	Repositories: [],
	// }
	// if err != nil {
	// 	log.Print(err)
	// 	rw.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	b, err := json.Marshal(nil)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(b))
}
