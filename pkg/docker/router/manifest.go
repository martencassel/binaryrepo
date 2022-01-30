package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog/log"
)

/*
	Check for a manifest
	HEAD /v2/<name>/manifests/<digest>
*/
func (router *DockerRouter) HasManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerreouter.HeadManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.HasManifest(rw, req)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.HasManifest(rw, req)
	}
}

/*
	Get a manifest
	GET /v2/<name>/manifests/<digest>
*/
func (router *DockerRouter) GetManifestHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.GetManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.GetManifestHandler(rw, req)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.GetManifestHandler(rw, req)
	}
}
