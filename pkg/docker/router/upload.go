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
func (router *DockerRouter) StartUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.StartUpload %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	urlParams := req.URL.Query()
	mountSha := urlParams.Get("mount")
	log.Info().Msgf("Repo-name: %s, mountSha: %s", repoName, mountSha)
	_repo := router.index.FindRepo(repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Type == repo.Local && _repo.PkgType == repo.Docker {
		router.registry.StartUpload(rw, req)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

/*
	Upload progress
	GET /v2/<name>/blobs/uploads/<uuid>
*/
func (router *DockerRouter) UploadProgress(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.UploadProgress %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Monolithic upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (router *DockerRouter) MonolithicUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.MonolithicUpload %s %s", req.Method, req.URL.Path)
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
	PATCH /v2/<name>/blobs/uploads/<uuid>
	Content-Length: <size of chunk>
	Content-Range: <start of range>-<end of range>
	Content-Type: application/octet-stream
	<Layer Chunk Binary Data>

	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
	Host: <registry host>
	Authorization: <scheme> <token>
	Content-Length: <length of data>
	Content-Type: application/octet-stream
	<binary data>

*/
func (router *DockerRouter) UploadChunk(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.UploadChunk %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Completed upload
	PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
*/
func (registry *DockerRouter) CompleteUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.CompleteUpload %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}

/*
	Cancel upload
	DELETE /v2/<name>/blobs/uploads/<uuid>
*/
func (registry *DockerRouter) CancelUpload(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("dockerrouter.CancelUpload %s %s", req.Method, req.URL.Path)
	log.Info().Msg("no implemented")
}
