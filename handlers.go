package binaryrepo

import (
	"net/http"
)

// This interface is a http handler that routes requests to local, remote or group docker registry handlers.
type DockerRegistryRouter interface {

	// HEAD /repo/{repo-name}/v2/<name>/manifests/<reference>
	HasManifest(rw http.ResponseWriter, req *http.Request)

	// GET /repo/{repo-name}/v2/<name>/manifests/<digest>
	GetManifestHandler(rw http.ResponseWriter, req *http.Request)

	// GET /repo/{repo-name}/v2/<name>/blobs/<digest>
	DownloadLayer(rw http.ResponseWriter, req *http.Request)

	// HEAD /repo/{repo-name}/v2/<name>/blobs/<digest>
	HasLayer(rw http.ResponseWriter, req *http.Request)

	// DELETE /repo/{repo-name}/v2/<name>/blobs/<digest>
	DeleteLayer(rw http.ResponseWriter, req *http.Request)

	// POST /repo/{repo-name}/v2/<name>/blobs/uploads/
	StartUpload(rw http.ResponseWriter, req *http.Request)

	// GET /repo/{repo-name}/v2/<name>/blobs/uploads/<uuid>
	UploadProgress(rw http.ResponseWriter, req *http.Request)

	// PUT /repo/{repo-name}/v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	MonolithicUpload(rw http.ResponseWriter, req *http.Request)

	// PATCH|PUT /repo/{repo-name}/v2/<name>/blobs/uploads/<uuid>(?digest=<digest>)
	UploadChunk(rw http.ResponseWriter, req *http.Request)

	// PUT /repo/{repo-name}/v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	CompleteUpload(rw http.ResponseWriter, req *http.Request)

	// DELETE /repo/{repo-name}/v2/<name>/blobs/uploads/<uuid>
	CancelUpload(rw http.ResponseWriter, req *http.Request)

	// GET /repo/{repo-name}/v2/version
	VersionHandler(rw http.ResponseWriter, req *http.Request)
}


// This interface is a http handler that implements a local docker registry http handlers
type LocalDockerRegistryHandler interface {
	// HEAD /v2/<name>/blobs/<digest>
	HasLayer(rw http.ResponseWriter, req *http.Request);

	// GET /v2/<name>/blobs/<digest>
	DownloadLayer(rw http.ResponseWriter, req *http.Request)

	// DELETE /v2/<name>/blobs/<digest>
	DeleteLayer(rw http.ResponseWriter, req *http.Request)

	// HEAD /v2/<name>/manifests/<reference>
	HasManifest(rw http.ResponseWriter, req *http.Request)

	// GET /v2/<name>/manifests/<reference>
	GetManifestHandler(rw http.ResponseWriter, req *http.Request)

	// PUT /v2/<name>/manifests/<reference>
	PutManifest(rw http.ResponseWriter, req *http.Request)

	// POST /v2/<name>/blobs/uploads/
	StartUpload(rw http.ResponseWriter, req *http.Request)

	// GET /v2/<name>/blobs/uploads/<uuid>
	UploadProgress(rw http.ResponseWriter, req *http.Request)

	// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	MonolithicUpload(rw http.ResponseWriter, req *http.Request)

	// PATCH /v2/<name>/blobs/uploads/<uuid>
	UploadChunk(rw http.ResponseWriter, req *http.Request)

	// PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	CompleteUpload(rw http.ResponseWriter, req *http.Request)

	// DELETE /v2/<name>/blobs/uploads/<uuid>
	CancelUpload(rw http.ResponseWriter, req *http.Request)
}


// This interface is a http handler that implements a remote caching proxy for docker registry http handlers
type RemoteDockerRegistryHandler interface {

	// HEAD /v2/version
	VersionHandler(rw http.ResponseWriter, req *http.Request)

	// HEAD /v2/<name>/manifests/<digest>
	HasManifest(rw http.ResponseWriter, req *http.Request)

	// GET /v2/<name>/manifests/<reference>
	GetManifestHandler(rw http.ResponseWriter, req *http.Request)

	// GET /v2/<name>/blobs/<digest>
	DownloadLayer(rw http.ResponseWriter, req *http.Request)

	// PUT /v2/<name>/blobs/<digest>
	LayerPut(rw http.ResponseWriter, req *http.Request)

	// HEAD /v2/<name>/blobs/<digest>
	HasLayer(rw http.ResponseWriter, req *http.Request)

}


// This interface is a http handler that implements a group handler for group repositories
type DockerGroupHandler interface {
}