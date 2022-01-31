package registry

import (
	"io/ioutil"
	"log"
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
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		b, err := os.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631.json")
		if err != nil {
			t.Fatal(err)
		}
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(digest)
		index.AddRepo(repo.Repo{ID: 1, Name: "test-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "test-local",
			"name":      "test",
			"digest":    "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631",
		}
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		req = mux.SetURLVars(req, vars)
		registry.DownloadLayer(res, req)
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "7700", res.Header().Get("Content-Length"))
		assert.Equal(t, "7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631", res.Header().Get("Docker-Content-Digest"))
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
	})
	// HEAD /repo/{repo-name}/v2/<name>/blobs/<digest>
	t.Run("Check if blob exists", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		b, err := os.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631.json")
		if err != nil {
			t.Fatal(err)
		}
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(digest)
		index.AddRepo(repo.Repo{ID: 1, Name: "test-local", Type: repo.Local, PkgType: repo.Docker})
		registry := NewDockerRegistry(fs, index)
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "test-local",
			"name":      "test",
			"digest":    "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631",
		}
		req = mux.SetURLVars(req, vars)
		registry.HasLayer(res, req)
		// Assert
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, "7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631", res.Header().Get("Docker-Content-Digest"))
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
	})
	// DELETE /repo/{repo-name}/v2/<name>/blobs/<digest>
	t.Run("Deleting a Layer", func(t *testing.T) {
		// Arrange
		os.RemoveAll("/tmp/filestore")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		b, err := os.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631.json")
		if err != nil {
			t.Fatal(err)
		}
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(digest)
		index.AddRepo(repo.Repo{ID: 1, Name: "test-local", Type: repo.Local, PkgType: repo.Docker})
		// Act
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "test-local",
			"name":      "test",
			"digest":    "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631",
		}
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		req = mux.SetURLVars(req, vars)
		registry := NewDockerRegistry(fs, index)
		registry.DeleteLayer(res, req)
		// Assert
		assert.Equal(t, http.StatusAccepted, res.Code)
		assert.Equal(t, "0", res.Header().Get("Content-Length"))
		assert.Equal(t, "7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631", res.Header().Get("Docker-Content-Digest"))
	})
}
