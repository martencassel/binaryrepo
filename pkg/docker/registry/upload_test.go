package registry

import (
	"bytes"
	_ "crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestUploads(t *testing.T) {
	/*
		Initiate Resumable Upload
		POST /v2/<name>/blobs/uploads/?digest=<digest>
		Host: <registry host>
		Authorization: <scheme> <token>
		Content-Length: <length of blob>
		Content-Type: application/octect-stream
		<binary data>
	*/
	// POST /repo/{repo-name}/v2/<name>/blobs/upload/
	t.Run("Initiate a resumable blob upload with an empty request body.", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)

		// Act
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "", nil)
		req.Header.Add("Content-Length", "0")
		vars := map[string]string{
			"repo-name": "redis-local",
			"name":      "redis",
		}
		req = mux.SetURLVars(req, vars)

		// Assert
		registry.StartUpload(res, req)
		assert.Equal(t, http.StatusAccepted, res.Code)
		assert.Contains(t, res.Header().Get("Content-Length"), "0")
		assert.Contains(t, res.Header().Get("docker-distribution-api-version"), "registry/2.0")
		assert.True(t, IsValidUUID(res.Header().Get("Docker-Upload-UUID")))
		assert.Equal(t, res.Header().Get("Range"), "0-0")
		assert.Equal(t, res.Header().Get("connection"), "close")
		uploadPath := fmt.Sprintf("/tmp/filestore/uploads/%s", res.Header().Get("Docker-Upload-UUID"))
		assert.True(t, fileExists(uploadPath))
		assert.Equal(t, res.Header().Get("Location"), fmt.Sprintf("/v2/%s/blobs/uploads/%s", vars["repo-name"], res.Header().Get("Docker-Upload-UUID")))
	})
	/*
		Complete the upload in a single request

		POST /v2/<name>/blobs/uploads/?digest=<digest>
		Host: <registry host>
		Authorization: <scheme> <token>
		Content-Length: <length of blob>
		Content-Type: application/octect-stream
		<binary data>
	*/
	// POST /repo/{repo-name}/v2/<name>/blobs/upload/?digest=<digest>
	t.Run("Complete upload in a single request", func(t *testing.T) {
		// Arrange
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		b, err := os.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631.json")
		if err != nil {
			t.Fatal(err)
		}
		id := digest.FromBytes(b)

		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "test-local",
			"name":      "test",
			"reference": "latest",
		}
		req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(b))
		req = mux.SetURLVars(req, vars)
		req.Header.Add("Content-Length", fmt.Sprintf("%d", len(b)))
		req.Header.Add("Content-Type", "application/octet-stream")
		q := req.URL.Query()
		q.Add("digest", id.String())
		registry.UploadChunk(res, req)

		// Assert
		offset := len(b) - 1
		assert.Equal(t, http.StatusAccepted, res.Code)
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
		assert.Contains(t, "0", res.Header().Get("Content-Length"))
		assert.True(t, IsValidUUID(res.Header().Get("docker-upload-uuid")))
		assert.Equal(t, fmt.Sprintf("0-%d", offset), res.Header().Get("range"))
	})
	/*
		Get Upload Status

		GET /v2/<name>/blobs/uploads/<uuid>
		Host: <registry host>
	*/
	// GET /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Get Blob Upload Status", func(t *testing.T) {
		// Arrange
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)

		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "test-local",
			"name":      "test",
			"reference": "latest",
		}
		req, _ := http.NewRequest(http.MethodGet, "", nil)
		req = mux.SetURLVars(req, vars)
		registry.UploadProgress(res, req)

		// Assert
	})
	/*
		Chunked Upload

		PATCH /v2/<name>/blobs/uploads/<uuid>
		Content-Length: <size of chunk>
		Content-Range: <start of range>-<end of range>
		Content-Type: application/octet-stream

	*/
	// PATCH /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Upload a chunk of data to specified upload", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		uuid := uuid.New().String()
		uploadPath := fmt.Sprintf("/tmp/filestore/uploads/%s", uuid)
		err := ioutil.WriteFile(uploadPath, []byte(""), 0644)
		if err != nil {
			t.Fatal(err)
		}
		blob, err := ioutil.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631.json")
		if err != nil {
			t.Fatal(err)
		}
		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "redis-local",
			"name":      "redis",
			"uuid":      uuid,
		}
		body := bytes.NewBuffer(blob)
		req, _ := http.NewRequest(http.MethodPut, "", body)
		req = mux.SetURLVars(req, vars)
		req.Header.Add("Transfer-Encoding", "chunked")
		req.Header.Add("Accept-Encoding", "gzip")
		registry.UploadChunk(res, req)

		// Assert
		offset := len(blob) - 1
		assert.Equal(t, http.StatusAccepted, res.Code)
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
		assert.Contains(t, "0", res.Header().Get("Content-Length"))
		assert.True(t, IsValidUUID(res.Header().Get("docker-upload-uuid")))
		assert.Equal(t, fmt.Sprintf("0-%d", offset), res.Header().Get("range"))
		assert.Equal(t, "close", res.Header().Get("Connection"))
	})
	/*
		Monolithic Upload

		PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
		Content-Length: <size of layer>
		Content-Type: application/octet-stream
		<Layer Binary Data>
	*/
	t.Run("Monolithic Upload", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	/*
		Complete Upload

		PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>
		Content-Length: <size of chunk>
		Content-Range: <start of range>-<end of range>
		Content-Type: application/octet-stream

		<Last Layer Chunk Binary Data>
	*/
	// PUT /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Complete the upload", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	/*
		Canceling an Upload
	*/
	// DELETE /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Cancel upload", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
}
