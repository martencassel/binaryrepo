package dockerproxy

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
)

// PathGetManifest URL.
const PathGetManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathGetManifest2 = "/repo/{repo-name}/v2/{namespace}/{namespace2}/manifests/{reference}"

// GetManifestHandler implements GET baseURL/repo/v2/namespace/manifests/reference
func (p *DockerProxyApp) GetManifestHandler(w http.ResponseWriter, req *http.Request) {
	opt := GetOptions(req)
	log.Printf("%s %s\n", req.Method, req.URL.Path)
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// PathHeadManifest URL.
const PathHeadManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathHeadManifest2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/manifests/{reference}"

func (p *DockerProxyApp) HeadManifestHandler(w http.ResponseWriter, req *http.Request) {
	opt := GetOptions(req)
	var _repo *repo.Repo
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo = p.index.FindRepo(repoName)
	if _repo == nil {
		log.Printf("Repo %s was not found", repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reference := vars["reference"]
	log.Print(reference, opt)
}
