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

// PUT /repo/docker-local/v2/test/blobs/uploads/8410da2d-9107-11ec-b583-0242ac1f0003

// Push endpoints
const PathUploadBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/"
const PathUploadBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/uploads/"
const PathUploadBlob3 = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/{uuid}"
//Path:"/repo/docker-local/repo/docker-local/v2/docker-local/blobs/uploads/84fa28f8-90fd-11ec-a736-0242ac1e0003"

const PathUploadBlob4 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/uploads/{uuid}"

const PathPutManifest = "/repo/{repo-name}/v2/{name}/manifests/{reference}"

func (router *DockerRouter) RegisterHandlers(r *mux.Router) {
	log.Info().Msg("Registering handlers")
	// GET /v2/
	r.HandleFunc(VerionHandlerPath1, router.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(VerionHandlerPath2, router.VersionHandler).Methods(http.MethodGet)
	// GET /v2/<name>/blobs/<digest>
	// POST /v2/<name>/blobs/uploads
	// HEAD /v2/<name>/blobs/<digest>
	// GET  /v2/<name>/blobs/uploads/<uuid>
	// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	// PATCH /v2/<name>/blobs/uploads/<uuid>
	// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	// DELETE /v2/<name>/blobs/uploads/<uuid>
	// DELETE /v2/<name>/blobs/<digest>
	// PUT /v2/<name>/manifests/<reference>
	// HEAD /v2/<name>/manifests/<reference>
	// GET /v2/<name>/manifests/<reference>
	// GET /v2/<name>/blobs/<digest>
	r.HandleFunc(dockerproxy.PathHeadManifest1, router.HasManifest).Methods(http.MethodHead)
	r.HandleFunc(dockerproxy.PathGetManifest2, router.HasManifest).Methods(http.MethodHead)

	r.HandleFunc(PathGetManifest1, router.GetManifestHandler).Methods(http.MethodGet)
	r.HandleFunc(PathGetManifest2, router.GetManifestHandler).Methods(http.MethodGet)

	r.HandleFunc(PathPutManifest, router.registry.PutManifest).Methods(http.MethodPut)

	r.HandleFunc(PathGetBlob1, router.DownloadLayer).Methods(http.MethodGet)
	r.HandleFunc(PathGetBlob2, router.DownloadLayer).Methods(http.MethodGet)


	r.HandleFunc(PathGetBlob1, router.HasLayer).Methods(http.MethodHead)
	r.HandleFunc(PathGetBlob2, router.HasLayer).Methods(http.MethodHead)

	// Starting an upload. HTTP POST
	r.HandleFunc(PathUploadBlob1, router.StartUpload).Methods(http.MethodPost)
	r.HandleFunc(PathUploadBlob2, router.StartUpload).Methods(http.MethodPost)

	// Chunked upload. HTTP PATCH
	r.HandleFunc(PathUploadBlob3, router.UploadChunk).Methods(http.MethodPatch)
	r.HandleFunc(PathUploadBlob4, router.UploadChunk).Methods(http.MethodPatch)

	// Monolithic upload or Complete the Upload. HTTP PUT
	r.HandleFunc(PathUploadBlob3, router.UploadChunk).Methods(http.MethodPut)
	r.HandleFunc(PathUploadBlob4, router.UploadChunk).Methods(http.MethodPut)

	r.HandleFunc(PathGetBlob1, router.DeleteLayer).Methods(http.MethodDelete)
	r.HandleFunc(PathGetBlob1, router.DeleteLayer).Methods(http.MethodDelete)


}


/*
	1. Monolithic Upload
	PUT	/v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	A monolithic upload is simply a chunked upload with a
	single chunk and may be favored by clients that would like
	to avoided the complexity of chunking.
	To carry out a “monolithic” upload, one can simply put
	the entire content blob to the provided URL:
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	Content-Length: <size of layer>
	Content-Type: application/octet-stream

	<Layer Binary Data>

	2. Completed Upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	Content-Length: <size of chunk>
	Content-Range: <start of range>-<end of range>
	Content-Type: application/octet-stream
	<Last Layer Chunk Binary Data>

	3. PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	Content-Length: 0
	Content-Type: application/octet-stream
	<Empty Body>

	201 Created
	Location: /v2/<name>/blobs/<digest>
	Content-Length: 0
	Docker-Content-Digest: <digest>
*/