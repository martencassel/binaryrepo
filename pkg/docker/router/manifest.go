package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog/log"
)

func (router *DockerRouter) HeadManifestHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerreouter.HeadManifest %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.HeadManifestHandler(w, r)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.HeadManifestHandler(w, r)
	}
}

func (router *DockerRouter) GetManifestHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.GetManifest %s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.GetManifestHandler(w, r)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.GetManifestHandler(w, r)
	}
}
