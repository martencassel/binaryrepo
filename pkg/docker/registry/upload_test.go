package dockerregistry

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestUploads(t *testing.T) {
	// Initiate Blob Upload.
	//
	// HTTP POST /v2/<name>/blobs/uploads/?digest=<digest>
	t.Run("Start an upload", func(t *testing.T) {
		// Arrange
		uploadUuid := "7315aa70-db4e-11ec-bf5c-0242ac1e0003"
//		dgst := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		repoName := "docker-local"
		namespace := "alpine"
		location := fmt.Sprintf("/v2/%s/blobs/uploads/%s", namespace, uploadUuid)
		repo := &binaryrepo.Repo{ Name: repoName }
		fs, rs, tag, uploader := CreateMocks()
		rs.On("FindRepo", repoName).Return(repo)
		uid, _ := uuid.Parse(uploadUuid)
		uploader.On("CreateUpload").Return(uid, nil)
//		digest_ := digest.Digest(dgst)
//		fs.On("WriteFile", mock.Anything).Return(digest_, nil)
		registry := NewDockerRegistryHandler(rs, fs, tag, uploader)
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "", nil)
		req.Header.Add("Content-Length", "0")
		req.Header.Add("Accept-Encoding", "gzip")
		// q := req.URL.Query()
        // q.Add("digest", dgst)
        // req.URL.RawQuery = q.Encode()
		vars := map[string]string{
			"repo-name": repoName,
			"namespace": namespace,
		}
		req = mux.SetURLVars(req, vars)

		// Act
		registry.StartUpload(res, req)

		// Assert
		assert.Equal(t, "0", res.Header().Get("Content-Length"))
		assert.Equal(t, uploadUuid, res.Header().Get("Docker-Upload-Uuid"))
		assert.Equal(t, location, res.Header().Get("Location"))
		assert.Equal(t, "0-0", res.Header().Get("Range"))
		assert.Equal(t, http.StatusAccepted, res.Code)
	})

	t.Run("Upload a chunk", func(t *testing.T) {
		// Arrange
		repoName := "docker-local"
		namespace := "alpine"
		uploadUuid := "7315aa70-db4e-11ec-bf5c-0242ac1e0003"
		dgst := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		location := fmt.Sprintf("/v2/%s/blobs/uploads/%s", namespace, uploadUuid)
		fs, rs, tag, uploader := CreateMocks()
		uploader.On("Exists", uploadUuid).Return(true)
		b, err := os.ReadFile("./testdata/2408cc74d12b6cd092bb8b516ba7d5e290f485d3eb9672efc00f0583730179e8")
		assert.NoError(t, err)
		uploader.On("WriteFile", uploadUuid, b).Return(nil)
		digest_, _ := digest.Parse(dgst)
		fs.On("WriteFile", b).Return(digest_, nil)

		if err != nil {
			t.Fatal(err)
		}
		vars := map[string]string{
			"repo-name": repoName,
			"namespace": namespace,
			"uuid": 	uploadUuid,
		}
		res := httptest.NewRecorder()
		body := bytes.NewBuffer(b)
		req, _ := http.NewRequest(http.MethodPut, "", body)
		req.Header.Add("Content-Length", "0")
		req.Header.Add("Accept-Encoding", "gzip")
		req = mux.SetURLVars(req, vars)
		registry := NewDockerRegistryHandler(rs, fs, tag, uploader)

		// Act
		registry.ChunkedUpload(res, req)

		// Assert
		assert.Equal(t, "0", res.Header().Get("Content-Length"))
		assert.Equal(t, uploadUuid, res.Header().Get("Docker-Upload-Uuid"))
		assert.Equal(t, location, res.Header().Get("Location"))
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, fmt.Sprintf("0-%d", len(b)-1), res.Header().Get("Range"))
	})

	t.Run("Complete an upload", func(t *testing.T) {
		// Arrange
		repoName := "docker-local"
		namespace := "alpine"
		uploadUuid := "7315aa70-db4e-11ec-bf5c-0242ac1e0003"
		dgst := "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		b, err := os.ReadFile("./testdata/2408cc74d12b6cd092bb8b516ba7d5e290f485d3eb9672efc00f0583730179e8")
		assert.NoError(t, err)
		fs, rs, tag, uploader := CreateMocks()
		uploader.On("Exists", uploadUuid).Return(true)
		uploader.On("ReadUpload", uploadUuid).Return(b, nil)
		uploader.On("WriteFile", uploadUuid, b).Return(nil)
		digest_, _ := digest.Parse(dgst)
		fs.On("WriteFile", b).Return(digest_, nil)

		req, _ := http.NewRequest(http.MethodPut, "", nil)
		req.Header.Add("Content-Length", "0")
		req.Header.Add("Accept-Encoding", "gzip")
		q := req.URL.Query()
        q.Add("digest", dgst)
        req.URL.RawQuery = q.Encode()
		vars := map[string]string{
			"repo-name": repoName,
			"namespace": namespace,
			"uuid": 	uploadUuid,
		}
		res := httptest.NewRecorder()
		req = mux.SetURLVars(req, vars)
		registry := NewDockerRegistryHandler(rs, fs, tag, uploader)

		// Act
		registry.CompleteUpload(res, req)

		// Assert
		assert.Equal(t, http.StatusCreated, res.Code)
	})
}