package dockerproxy

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// PathVersion URL.
const PathVersionUrl1 = "/repo/{repo-name}/v2"
const PathVersionUrl2 = "/repo/{repo-name}/v2/"

// VersionHandler implements GET baseURL/repo/v2/
func (p *DockerProxyApp) VersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	if p.index.FindRepo(repoName) == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
