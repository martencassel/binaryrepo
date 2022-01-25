package docker

import (
	"net/http"

	"github.com/gorilla/mux"

	proxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
	registry "github.com/martencassel/binaryrepo/pkg/docker/registry"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

type DockerRouter struct {
	proxy    *proxy.DockerProxyApp
	registry *registry.DockerRegistry
	index    *repo.RepoIndex
}

func NewDockerRouter(proxy *proxy.DockerProxyApp, registry *registry.DockerRegistry, repoIndex *repo.RepoIndex) *DockerRouter {
	return &DockerRouter{
		proxy:    proxy,
		registry: registry,
		index:    repoIndex,
	}
}

const PathVersionUrl1 = "/repo/{repo-name}/v2"
const PathVersionUrl2 = "/repo/{repo-name}/v2/"
const PathHeadManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathHeadManifest2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/manifests/{reference}"
const PathGetManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathGetManifest2 = "/repo/{repo-name}/v2/{namespace}/{namespace2}/manifests/{reference}"
const PathGetBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/{digest}"
const PathGetBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/{digest}"

func (router *DockerRouter) RegisterHandlers(r *mux.Router) {
	r.HandleFunc(PathVersionUrl1, router.versionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathVersionUrl2, router.versionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathHeadManifest1, router.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathHeadManifest2, router.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathGetManifest1, router.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetManifest2, router.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob1, router.DownloadLayer).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob2, router.DownloadLayer).Methods(http.MethodGet)
}

func (router *DockerRouter) versionHandler(w http.ResponseWriter, r *http.Request) {
	////log.Info().Msgf("%s %s\n", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Remote && _repo.PkgType == repo.Docker {
		router.proxy.VersionHandler(w, r)
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.VersionHandler(w, r)
	}
}

func (router *DockerRouter) HeadManifestHandler(w http.ResponseWriter, r *http.Request) {
	////log.Info().Msgf("%s %s\n", r.Method, r.URL.Path)
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
	////log.Info().Msgf("%s %s\n", r.Method, r.URL.Path)
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

func (router *DockerRouter) DownloadLayer(w http.ResponseWriter, r *http.Request) {
	////log.Info().Msgf("%s %s\n", r.Method, r.URL.Path)
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
