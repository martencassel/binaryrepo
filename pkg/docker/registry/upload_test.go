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
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestUploads(t *testing.T) {
	// POST /repo/{repo-name}/v2/<name>/blobs/upload
	t.Run("Start an upload", func(t *testing.T) {
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
		registry.StartUpload(res, req)

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
	// GET /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Get Blob Upload Status", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	t.Run("Monolithic Upload", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	// PATCH /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Upload a chunk of data", func(t *testing.T) {
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
	// PUT /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Complete the upload", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	// DELETE /repo/{repo-name}/v2/<name>/blobs/upload/<uuid>
	t.Run("Cancel upload", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})

}
