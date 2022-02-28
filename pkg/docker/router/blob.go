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
func (router *DockerRouter) DownloadLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.DownloadLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.DownloadLayer(rw, req)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.DownloadLayer(rw, req)
	}
}

/*
	Check for layer
	HEAD /v2/<name>/blobs/<digest>
*/
func (router *DockerRouter) HasLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.HasLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.HasLayer(rw, req)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.HasLayer(rw, req)
	}
}

/*
	Deleting a Blob
	DELETE /v2/<name>/blobs/<digest>
*/
func (router *DockerRouter) DeleteLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("router.DeleteLayer %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.DeleteLayer(rw, req)
	}
	rw.WriteHeader(http.StatusNotFound)
}
