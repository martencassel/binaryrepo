package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/rs/zerolog/log"
)

/*
	Start an upload
	POST /v2/<name>/blobs/uploads
*/
func (registry *DockerRegistry) StartUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.StartUpload %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	namespace := vars["namespace"]
	repoName := vars["repo-name"]
	repo := registry.index.FindRepo(repoName)
	if repo == nil || namespace == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	uuid, _ := uuid.NewUUID()

	err := registry.uploader.WriteFile(uuid.String(), []byte{})
	if err != nil {
		log.Error().Msg(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	loc := fmt.Sprintf("/repo/%s/v2/%s/blobs/uploads/%s", repoName, namespace, uuid)
	rw.Header().Set("Content-Length", "0")
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Docker-Upload-UUID", uuid.String())
	rw.Header().Set("Range", "0-0")
	rw.Header().Set("Location", loc)
	rw.Header().Set("Connection", "close")
	rw.WriteHeader(http.StatusAccepted)
}

/*
	Upload progress
	GET /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRegistry) UploadProgress(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.UploadProgress %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Monolithic upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (registry *DockerRegistry) MonolithicUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.MonolithicUpload %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Chunked upload
	PATCH /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRegistry) UploadChunk(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.UploadChunk %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["namespace"]
	repoName := vars["repo-name"]
	uuid := vars["uuid"]
	log.Printf("%s /repo/%s/v2/%s/blobs/uploads/%s", req.Method, repoName, name, uuid)
	/*if req.Method == http.MethodPatch) || !(req.Method == http.MethodPut) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}*/
	if registry.index.FindRepo(repoName) == nil || name == "" || uuid == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if req.Body == http.NoBody || req.ContentLength == 0 {
		urlParams := req.URL.Query()
		digest := urlParams.Get("digest")

		b, err := registry.uploader.ReadUpload(uuid)
		if err != nil {
			log.Error().Msg(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info().Msgf("Size: %d", len(b))
		dgst, err := registry.fs.WriteFile(b)
		if err != nil {
			log.Error().Msg(err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info().Msgf("Digest: %s", dgst)
		log.Info().Msgf("Digest: %s", digest)
		rw.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/%s", repoName, digest))
		rw.Header().Set("Content-Length", "0")
		rw.WriteHeader(http.StatusCreated)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !IsValidUUID(uuid) {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.Header.Get("Range") != "" {
		rw.WriteHeader(http.StatusNotImplemented)
		return
	}
	if !registry.uploader.Exists(uuid) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	registry.uploader.AppendFile(uuid, body)

	location := fmt.Sprintf("/repo/%s/v2/%s/blobs/uploads/%s", repoName, name, uuid)
	offset := len(body) - 1
	rw.Header().Set("Content-Length", "0")
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Docker-Upload-UUID", uuid)
	rw.Header().Set("Location", location)
	rw.Header().Set("Range", fmt.Sprintf("0-%d", offset))
	//rw.Header().Set("connection", "close")
	rw.WriteHeader(http.StatusAccepted)
}

/*
	Completed upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (registry *DockerRegistry) CompleteUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.CompleteUpload %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Cancel upload
	DELETE /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRegistry) CancelUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.CancelUpload %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}
