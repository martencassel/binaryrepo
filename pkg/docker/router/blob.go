package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog/log"
)

/*
	Get blob
	GET /v2/<name>/blobs/<digest>
*/
func (router *DockerRouter) DownloadLayer(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.DownloadLayer %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.DownloadLayer(w, r)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.DownloadLayer(w, r)
	}
}

/*
	Check for layer
	HEAD /v2/<name>/blobs/<digest>
*/
func (router *DockerRouter) HasLayer(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.HasLayer %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.HasLayer(w, r)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.HasLayer(w, r)
	}
}
