package dockerregistry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/mocks"
	"github.com/stretchr/testify/assert"
)


func CreateMocks() (*mocks.FileStoreMock,
					*mocks.RepoStoreMock,
					*mocks.TagStoreMock,
					*mocks.UploaderMock) {
	fileStoreMock := mocks.NewFileStoreMock()
	repoStoreMock := mocks.NewRepoStoreMock()
	tagStoreMock := mocks.NewTagStoreMock()
	uploader := mocks.NewUploaderMock()
	return fileStoreMock, repoStoreMock, tagStoreMock, uploader
}

func TestPing(t *testing.T) {
	t.Run("Ping registry", func(t *testing.T) {
		// Arrange
		fs, rs, tag, uploader := CreateMocks()

		registry := NewDockerRegistryHandler(rs, fs, tag, uploader)

		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "", nil)
		req.Header.Add("Content-Length", "0")
		vars := map[string]string{
			"repo-name": "redis-local",
			"name":      "redis",
		}
		req = mux.SetURLVars(req, vars)

		// Act
		registry.VersionHandler(res, req)

		// Assert
		assert.Equal(t, "registry/2.0", res.Header().Get("Docker-Distribution-API-Version"))
		assert.Equal(t, http.StatusOK, res.Code)
	})
}