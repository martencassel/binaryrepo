package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog/log"
)

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

func (router *DockerRouter) initiateBlobUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.initiateBlobUpload %s %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	urlParams := r.URL.Query()
	mountSha := urlParams.Get("mount")
	log.Info().Msgf("Repo-name: %s, mountSha: %s", repoName, mountSha)
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.InitiateBlobUpload(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (router *DockerRouter) blobUploadChunk(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.blobUploadChunk %s %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	uuid := vars["uuid"]
	log.Info().Msgf("Repo-name: %s, uuid: %s\n", repoName, uuid)
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.UploadBlobChunk(w, r)
	}
	w.WriteHeader(http.StatusNotFound)
}
