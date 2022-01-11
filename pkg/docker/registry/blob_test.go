package registry

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

func TestBlob(t *testing.T) {
	t.Run("Initiate a blob upload", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)

		// Act
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "", nil)
		req.Header.Add("Accept-Encoding", "gzip")
		req.Header.Add("Content-Length", "0")
		vars := map[string]string{
			"repo-name": "redis-local",
			"name":      "redis",
		}
		req = mux.SetURLVars(req, vars)
		registry.InitBlobUpload(res, req)

		// Assert
		assert.Contains(t, res.Header().Get("Content-Length"), "0")
		assert.Contains(t, res.Header().Get("docker-distribution-api-version"), "registry/2.0")
		assert.True(t, IsValidUUID(res.Header().Get("Docker-Upload-UUID")))
		assert.Equal(t, res.Header().Get("Range"), "0-0")
		assert.Equal(t, res.Header().Get("connection"), "close")
		assert.Equal(t, http.StatusAccepted, res.Code)
		uploadPath := fmt.Sprintf("/tmp/filestore/uploads/%s", res.Header().Get("Docker-Upload-UUID"))
		assert.True(t, fileExists(uploadPath))
	})

	t.Run("Upload a chunk of data", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		uuid := uuid.New().String()
		uploadPath := fmt.Sprintf("/tmp/filestore/uploads/%s", uuid)
		ioutil.WriteFile(uploadPath, []byte(""), 0644)
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

	t.Run("Check for an existing layer", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "test-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "test-local",
			"name":      "test",
			"digest":    "sha256:43d70e856f23a163ea50230761ff1dbad23e127553c860ac3d6c677c712a059d",
		}
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		req = mux.SetURLVars(req, vars)
		registry.ExistsBlob(res, req)
		// Assert
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
	})
	t.Run("Put a manifest", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		index.AddRepo(repo.Repo{ID: 1, Name: "redis-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "redis-local",
			"name":      "redis",
			"reference": "latest",
		}
		blob, err := ioutil.ReadFile("./testdata/redis-manifest.json")
		if err != nil {
			t.Fatal(err)
		}
		body := bytes.NewBuffer(blob)
		req, _ := http.NewRequest(http.MethodPut, "", body)
		req = mux.SetURLVars(req, vars)
		req.Header.Add("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		req.Header.Add("Accept-Encoding", "gzip")
		registry.PutManifest(res, req)
		// Assert
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, "sha256:563888f63149e3959860264a1202ef9a644f44ed6c24d5c7392f9e2262bd3553", res.Header().Get("docker-content-digest"))
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
		assert.Contains(t, "redis/manifests/sha256:563888f63149e3959860264a1202ef9a644f44ed6c24d5c7392f9e2262bd3553", res.Header().Get("Location"))
		assert.Contains(t, "0", res.Header().Get("Content-Length"))
		assert.Equal(t, "close", res.Header().Get("Connection"))
	})
}
