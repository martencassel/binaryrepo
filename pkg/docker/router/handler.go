package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	dockerproxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
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

// Pull endpoints
const PathGetManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathGetManifest2 = "/repo/{repo-name}/v2/{namespace}/{namespace2}/manifests/{reference}"
const PathGetBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/{digest}"
const PathGetBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/{digest}"

// Push endpoints
const PathUploadBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/"
const PathUploadBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/uploads/"
const PathUploadBlob3 = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/{uuid}"
const PathUploadBlob4 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/uploads/{uuid}"

func (router *DockerRouter) RegisterHandlers(r *mux.Router) {
	log.Info().Msg("Registering handlers")
	r.HandleFunc(PathVersionUrl1, router.versionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathVersionUrl2, router.versionHandler).Methods(http.MethodGet)
	r.HandleFunc(dockerproxy.PathHeadManifest1, router.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(dockerproxy.PathGetManifest2, router.HeadManifestHandler).Methods(http.MethodHead)
	r.HandleFunc(PathGetManifest1, router.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetManifest2, router.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob1, router.DownloadLayer).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob2, router.DownloadLayer).Methods(http.MethodGet)
	// Push endpoints
	r.HandleFunc(PathUploadBlob1, router.initiateBlobUpload).Methods(http.MethodPost)
	r.HandleFunc(PathUploadBlob2, router.initiateBlobUpload).Methods(http.MethodPost)
	r.HandleFunc(PathUploadBlob3, router.blobUploadChunk).Methods(http.MethodPatch)
	r.HandleFunc(PathUploadBlob4, router.blobUploadChunk).Methods(http.MethodPatch)
	r.HandleFunc(PathUploadBlob3, router.blobUploadChunk).Methods(http.MethodPatch)
	r.HandleFunc(PathUploadBlob4, router.blobUploadChunk).Methods(http.MethodPatch)
}
