package dockerrouter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog/log"
)

/*
	Start an upload
	POST /v2/<name>/blobs/uploads
*/
func (router *DockerRouter) StartUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.StartUpload %s %s", r.Method, r.URL.Path)
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
		router.registry.StartUpload(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

/*
	Upload progress
	GET /v2/<name>/blobs/uploads/<uuid>
*/
func (router *DockerRouter) UploadProgress(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.UploadProgress %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Monolithic upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (router *DockerRouter) MonolithicUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.MonolithicUpload %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	The client can specify a range header and only include that part of the layer file.
	PATCH /v2/<name>/blobs/uploads/<uuid>
	Content-Length: <size of chunk>
	Content-Range: <start of range>-<end of range>
	Content-Type: application/octet-stream
	<Layer Chunk Binary Data>

	202 Accepted
	Location: /v2/<name>/blobs/uploads/<uuid>
	Range: bytes=0-<offset>
	Content-Length: 0
	Docker-Upload-UUID: <uuid>
*/
/*
	Upload a chunk of data for the specified found
	PATCH /v2/<name>/blobs/<uuid>
*/
func (router *DockerRouter) UploadChunk(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.UploadChunk %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Completed upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (registry *DockerRouter) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.CompleteUpload %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Cancel upload
	DELETE /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRouter) CancelUpload(w http.ResponseWriter, r *http.Request) {
	log.Info().Msgf("dockerrouter.CancelUpload %s %s", r.Method, r.URL.Path)
	log.Info().Msg("no implemented")
}
