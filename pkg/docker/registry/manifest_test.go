package registry

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestManifest(t *testing.T) {
	// GET /v2/<name>/manifests/<reference>
	t.Run("Pull a manifest", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	// HEAD /v2/<name>/manifests/<reference>
	t.Run("Check if manifest exists", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})
	// PUT /v2/<name>/manifests/<reference>
	t.Run("Push a manifest", func(t *testing.T) {
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
		assert.Contains(t, "sha256:3b97f312f894b02e4572bf831bad6343b45f5a08280af34ee2001140f342fe72", res.Header().Get("docker-content-digest"))
		assert.Contains(t, "registry/2.0", res.Header().Get("docker-distribution-api-version"))
		assert.Contains(t, "/v2/redis/manifests/latest", res.Header().Get("Location"))
		//		assert.Contains(t, "0", res.Header().Get("Content-Length"))
	})
	// DELETE /v2/<name>/manifests/<reference>
	t.Run("Deleting an manifest", func(t *testing.T) {
		// Arrange
		// Act
		// Assert
	})

}
