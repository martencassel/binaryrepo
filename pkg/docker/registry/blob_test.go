package registry

import (
	"errors"
	"fmt"
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

func TestBlobUpload(t *testing.T) {
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
