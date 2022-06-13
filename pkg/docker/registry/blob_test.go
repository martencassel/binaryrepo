package dockerregistry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBlobs(t *testing.T) {

	t.Run("Check if blob exists", func(t *testing.T) {
		// Arrange
		repoName := "docker-local"
		namespace := "alpine"
		base_digest := "0c4164589f761881ef975352441e796636759369b39de5fa49bae37126035715"
		dgst := "sha256:" + base_digest

		digest_, err := digest.Parse(dgst)
		assert.NoError(t, err)

		// Filestore, Exists
		fs, rs, tag, uploader := CreateMocks()
		fs.On("Exists", digest_).Return(true)

		// Filestore, ReadFile
		b, err := os.ReadFile("./testdata/2408cc74d12b6cd092bb8b516ba7d5e290f485d3eb9672efc00f0583730179e8")
		assert.NoError(t, err)
		fs.On("ReadFile", digest_).Return(b, nil)

		// Request
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		req.Header.Add("Content-Length", "0")
		req.Header.Add("Accept-Encoding", "gzip")

		vars := map[string]string{
			"repo-name": repoName,
			"namespace": namespace,
			"digest": 	dgst,
		}
		res := httptest.NewRecorder()
		req = mux.SetURLVars(req, vars)

		registry := NewDockerRegistryHandler(rs, fs, tag, uploader)

		// Act
		registry.ExistingLayer(res, req)

		// Assert
		assert.Equal(t, "application/octet-stream", res.Header().Get("Content-Type"))
		assert.Equal(t, "8060", res.Header().Get("Content-Length"))
		assert.Equal(t, base_digest, res.Header().Get("Docker-Content-Digest"))
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Get a blob", func(t *testing.T) {
		// Arrange
		repoName := "docker-local"
		namespace := "alpine"
		base_digest := "2408cc74d12b6cd092bb8b516ba7d5e290f485d3eb9672efc00f0583730179e8"
		dgst := "sha256:" + base_digest
		digest_, err := digest.Parse(dgst)
		assert.NoError(t, err)

		// Mocks
		fs, rs, _, _ := CreateMocks()
		rs.On("Exists", mock.Anything, repoName).Return(true, nil)
		b, err := os.ReadFile(fmt.Sprintf("./testdata/%s", base_digest))
		assert.NoError(t, err)
		fs.On("Exists", digest_).Return(true)
		fs.On("ReadFile", digest_).Return(b, nil)

		// Request, Response
		vars := map[string]string{
			"repo-name": repoName,
			"namespace": namespace,
			"digest": 	dgst,
		}
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		res := httptest.NewRecorder()
		req = mux.SetURLVars(req, vars)

		// Handler
		registry := NewDockerRegistryHandler(rs, fs, nil, nil)

		// Act
		registry.GetLayer(res, req)

		// Assert
		assert.Equal(t, "application/octet-stream", res.Header().Get("Content-Type"))
		assert.Equal(t, "8060", res.Header().Get("Content-Length"))
		assert.Equal(t, base_digest, res.Header().Get("Docker-Content-Digest"))
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Delete a blob if it exists", func(t *testing.T) {
		// Arrange
		dgst := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		digest_, err := digest.Parse(dgst)
		assert.NoError(t, err)
		vars := map[string]string{
			"repo-name": "docker-local",
			"namespace": "alpine",
			"digest": 	dgst,
		}
		// Request, Response
		req, _ := http.NewRequest(http.MethodHead, "", nil)
		res := httptest.NewRecorder()
		req = mux.SetURLVars(req, vars)

		// Mocks
		fs, rs, _, _ := CreateMocks()
		rs.On("Exists", mock.Anything, "docker-local").Return(true, nil)
		fs.On("Exists", digest_).Return(true)
		// Handler
		registry := NewDockerRegistryHandler(rs, fs, nil, nil)

		// Act
		registry.DeleteLayer(res, req)

		// Assert
		fs.On("Exists", digest_).Return(false)
		assert.Equal(t, http.StatusOK, res.Code)
	})
}
