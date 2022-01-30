package registry

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

func TestBlob(t *testing.T) {
	// GET /repo/{repo-name}/v2/<name>/blobs/<digest>
	t.Run("Pull a Layer", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	// HEAD /repo/{repo-name}/v2/<name>/blobs/<digest>
	t.Run("Check if blob exists", func(t *testing.T) {
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
	// DELETE /repo/{repo-name}/v2/<name>/blobs/<digest>
	t.Run("Deleting a Layer", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
}
