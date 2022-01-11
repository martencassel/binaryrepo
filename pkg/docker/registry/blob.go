package registry

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	digest "github.com/opencontainers/go-digest"
)

// PathInitBlobUpload URL.
const PathInitBlobUpload = "/repo/{repo-name}/v2/{name}/blobs/upload"

// Initiate a resumable blob upload.
func (registry *DockerRegistry) InitBlobUpload(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	if registry.index.FindRepo(repoName) == nil || name == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("%s /v2/%s/blobs/uploads", http.MethodPost, name)
	uuid, _ := uuid.NewUUID()
	err := ioutil.WriteFile(fmt.Sprintf("%s/uploads/%s", registry.fs.BasePath, uuid.String()), []byte{}, 0644)
	if err != nil {
		log.Fatal(err)
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

// PathChunkedUpload URL.
const PathChunkedUpload = "/repo/{repo-name}/v2/{name}/blobs/uploads/{uuid}"

// Upload a chunk of data for the specified upload.
func (registry *DockerRegistry) UploadChunk(w http.ResponseWriter, r *http.Request) {
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

// PathExistsBlob URL.
const PathExistsBlob = "/repo/{repo-name}/v2/{name}/blobs/{uuid}"

func (registry *DockerRegistry) ExistsBlob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	repoName := vars["repo-name"]
	d := vars["digest"]
	if r.Method != http.MethodHead {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if registry.index.FindRepo(repoName) == nil || name == "" || d == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(d)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !registry.fs.Exists(digest) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := registry.fs.ReadFile(digest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("docker-distribution-api-version", "registry/2.0")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("connection", "close")
	w.WriteHeader(http.StatusAccepted)
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
