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
	log.Printf("VersionHandler: %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	log.Printf("Repo Name: %s\n", repoName)
	repo := p.index.FindRepo(repoName)
	log.Println(r)
	if repo == nil {
		log.Printf("Repo %s was not found", repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
