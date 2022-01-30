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
func (registry *DockerRegistry) StartUpload(rw http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.StartUpload %s %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	if registry.index.FindRepo(repoName) == nil || name == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	uuid, _ := uuid.NewUUID()
	err := ioutil.WriteFile(fmt.Sprintf("%s/uploads/%s", registry.fs.BasePath, uuid.String()), []byte{}, 0644)
	if err != nil {
		log.Error().Msg(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	loc := fmt.Sprintf("/repo/%s/v2/%s/blobs/uploads/%s", repoName, name, uuid)
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
func (registry *DockerRegistry) UploadProgress(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.UploadProgress %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Monolithic upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (registry *DockerRegistry) MonolithicUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.MonolithicUpload %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Chunked upload
	PATCH /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRegistry) UploadChunk(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.UploadChunk %s %s", r.Method, r.URL.Path)
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	uuid := vars["uuid"]
	log.Printf("%s /repo/%s/v2/%s/blobs/uploads/%s", r.Method, repoName, name, uuid)
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if registry.index.FindRepo(repoName) == nil || name == "" || uuid == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if r.Body == http.NoBody || r.ContentLength == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !IsValidUUID(uuid) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Range") != "" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	uploadPath := fmt.Sprintf("%s/uploads/%s", registry.fs.BasePath, uuid)
	if !fileExists(uploadPath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	location := fmt.Sprintf("/repo/%s/v2/%s/blobs/uploads/%s", repoName, name, uuid)
	offset := len(body) - 1
	w.Header().Set("Content-Length", "0")
	w.Header().Set("docker-distribution-api-version", "registry/2.0")
	w.Header().Set("Docker-Upload-UUID", uuid)
	w.Header().Set("Location", location)
	w.Header().Set("Range", fmt.Sprintf("0-%d", offset))
	w.Header().Set("connection", "close")
	w.WriteHeader(http.StatusAccepted)
}

/*
	Completed upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (registry *DockerRegistry) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.CompleteUpload %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Cancel upload
	DELETE /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRegistry) CancelUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("registry.CancelUpload %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}
