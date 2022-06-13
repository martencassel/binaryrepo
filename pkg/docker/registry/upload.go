package dockerregistry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opencontainers/go-digest"
	"github.com/rs/zerolog/log"
)

// Start an upload session.
// POST /v2/<name>/blobs/uploads/
func (registry *DockerRegistryHandler) StartUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.StartUpload %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	namespace := vars["namespace"]
	repoName := vars["repo-name"]
	urlParams := req.URL.Query()
	param_digest := urlParams.Get("digest")

	uuid, err := registry.uploader.CreateUpload()
	if err != nil {
		log.Error().Msg(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// repo := registry.index.FindRepo(repoName)
	// if repo == nil || namespace == "" {
	// 	rw.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// Check if monolithic upload. digest != ""

	// Non resumable upload
	if param_digest != "" {
		dgst, err := digest.Parse(param_digest)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info().Msgf("Digest: %s", dgst)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		dgst, err = registry.fs.WriteFile(body)
		if err != nil {
			log.Error().Msg(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info().Msgf("Digest: %s", dgst)
		loc := fmt.Sprintf("/repo/%s/v2/%s/blobs/%s", repoName, namespace, dgst)
		rw.Header().Set("Location", loc)
		rw.Header().Set("Content-Length", "0")
		rw.Header().Set("Docker-Upload-UUID", param_digest)
		rw.WriteHeader(http.StatusCreated)
		return
	}

	// err := registry.uploader.WriteFile(uuid.String(), []byte{})
	// if err != nil {
	// 	log.Error().Msg(err.Error())
	// 	rw.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	loc := fmt.Sprintf("/v2/%s/blobs/uploads/%s", namespace, uuid)
	rw.Header().Set("Content-Length", "0")
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Docker-Upload-UUID", uuid.String())
	rw.Header().Set("Range", "0-0")
	rw.Header().Set("Location", loc)
	rw.Header().Set("Connection", "close")
	rw.WriteHeader(http.StatusAccepted)
}

// Upload status
// GET /v2/<name>/blobs/uploads/<uuid>
func (registry *DockerRegistryHandler) UploadStatus(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.UploadStatus %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

// Chunked upload
// PATCH /v2/<name>/blobs/uploads/<uuid>
func (registry *DockerRegistryHandler) ChunkedUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.ChunkedUpload %s %s", req.Method, req.URL.Path)

	vars := mux.Vars(req)
	namespace := vars["namespace"]
	repoName := vars["repo-name"]
	uuid := vars["uuid"]
	log.Info().Msgf("namespace: %s, repo: %s, uuid: %s", namespace, repoName, uuid)

	if req.Body == http.NoBody || req.ContentLength == 0 {
		log.Info().Msg("no body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	exists := registry.uploader.Exists(uuid)
	if exists == false {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	err = registry.uploader.WriteFile(uuid, body)
	if err != nil {
		log.Error().Msg(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
	dgst, err := registry.fs.WriteFile(body)
	if err != nil {
		log.Error().Msg(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info().Msgf("Write file to filestore with Digest: %s", dgst)

	contentLength := req.Header.Get("Content-Length")
	contentRange := req.Header.Get("Content-Range")
	contentType := req.Header.Get("Content-Type")
	log.Info().Msgf("Content-Length: %s, Content-Range: %s, Content-Type: %s", contentLength, contentRange, contentType)

	location := fmt.Sprintf("/v2/%s/blobs/uploads/%s", namespace, uuid)
	offset := len(body) - 1
	rw.Header().Set("Content-Length", "0")
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Docker-Upload-UUID", uuid)
	rw.Header().Set("Location", location)
	rw.Header().Set("Range", fmt.Sprintf("0-%d", offset))

	if req.Method == http.MethodPut {
		rw.WriteHeader(http.StatusCreated)
	} else {
		rw.WriteHeader(http.StatusAccepted)
	}
}

// Complete status
// PUT /v2/<name>/blobs/uploads/<uuid>
func (registry *DockerRegistryHandler) CompleteUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.CompleteUpload %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	urlParams := req.URL.Query()
	namespace := vars["namespace"]
	repoName := vars["repo-name"]
	uuid := vars["uuid"]
	digest := urlParams.Get("digest")
	log.Info().Msgf("namespace: %s, repo: %s, uuid: %s, digest: %s\n", namespace, repoName, uuid, digest)
	if registry.uploader.Exists(uuid) == false {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	// Optionally, if all chunks have been uploaded, a PUT request with a digest parameter and a zero-length body
	// may be sent to complete and validate the upload.
	if req.Body == http.NoBody || req.ContentLength == 0 {

		b, err := registry.uploader.ReadUpload(uuid)
		if err != nil {
			log.Error().Msg(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		dgst, err := registry.fs.WriteFile(b)
		if err != nil {
			log.Error().Msg(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		location := fmt.Sprintf("/repo/%s/v2/%s/blobs/%s", repoName, namespace, digest)
		rw.Header().Set("Content-Length", "0")
		rw.Header().Set("docker-distribution-api-version", "registry/2.0")
		rw.Header().Set("Location", location)
		rw.Header().Set("Docker-Content-Digest", dgst.String())
		rw.WriteHeader(http.StatusCreated)
		return
	}
	rw.WriteHeader(http.StatusNotImplemented)
}

// Delete Blob upload
// DELETE /v2/<name>/blobs/uploads/<uuid>
func (registry *DockerRegistryHandler) DeleteBlobUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.DeleteBlobUpload %s %s", req.Method, req.URL.Path)
}