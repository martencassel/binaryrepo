package dockerregistry

import (
	"net/http"

	"github.com/gorilla/mux"
)

const PathCheckBlob        = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/{uuid}"
const PathStartUpload      = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/"
const PathUploadStatus     = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/{uuid}"
const PathChunkedUpload    = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/{uuid}"
const PathCompleteUpload   = "/repo/{repo-name}/v2/{namespace}/blobs/uploads/{uuid}"
const PathExistingLayer    = "/repo/{repo-name}/v2/{namespace}/blobs/{digest}"
const PathPushManifest     = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathExistingManifest = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathGetManifest      = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathListImageTags    = "/repo/{repo-name}/v2/{namespace}/tags/list"

func (registry *DockerRegistryHandler) RegisterHandlers(r *mux.Router) {

	r.HandleFunc(PathVersionUrl1, registry.VersionHandler).Methods(http.MethodGet)
	r.HandleFunc(PathVersionUrl2, registry.VersionHandler).Methods(http.MethodGet)

	r.HandleFunc(PathStartUpload, registry.StartUpload).Methods(http.MethodPost)
	r.HandleFunc(PathChunkedUpload, registry.ChunkedUpload).Methods(http.MethodPatch)
	r.HandleFunc(PathCompleteUpload, registry.CompleteUpload).Methods(http.MethodPut)
	r.HandleFunc(PathExistingLayer, registry.ExistingLayer).Methods(http.MethodHead)
	r.HandleFunc(PathExistingLayer, registry.GetLayer).Methods(http.MethodGet)
	r.HandleFunc(PathPushManifest, registry.UploadManifest).Methods(http.MethodPut)
	r.HandleFunc(PathListImageTags, registry.ListImageTags).Methods(http.MethodGet)

	r.HandleFunc(PathExistingManifest, registry.ExistingManifest).Methods(http.MethodHead)
	r.HandleFunc(PathGetManifest, registry.GetManifest).Methods(http.MethodGet)

}