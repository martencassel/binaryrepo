package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog/log"
)

const VerionHandlerPath1 = "/repo/{repo-name}/v2"
const VerionHandlerPath2 = "/repo/{repo-name}/v2/"

func (router *DockerRouter) VersionHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.versionHandler %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.VersionHandler(rw, req)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.VersionHandler(rw, req)
	}
}
