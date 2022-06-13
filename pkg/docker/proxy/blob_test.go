package dockerproxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/mocks"
	digest "github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func CreateMocks() (*mocks.FileStoreMock, *mocks.RepoStoreMock, *mocks.RegistryClientMock) {
	fileStoreMock := mocks.NewFileStoreMock()
	repoStoreMock := mocks.NewRepoStoreMock()
	registryClientMock := mocks.NewRegistryClientMock()
	return fileStoreMock, repoStoreMock, registryClientMock
}

func TestFetchBlob(t *testing.T) {
	t.Run("Fetch blob", func(t *testing.T) {
		d := "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631"
		digest, err := digest.Parse(d)
		if err != nil {
				t.Fatalf("Digest is invalid %s", err)
		}
		req, err := http.NewRequest("GET", "http://localhost:8081/repo/docker-remote/v2/redis/blobs/sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631", nil)
		if err != nil {
				t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		fileStoreMock, repoServiceMock, registryClientMock := CreateMocks()
		fileStoreMock.On("Exists", digest).Return(true)
		fileStoreMock.On("ReadFile", digest).Return([]byte("Hello World"), nil)
		r := &binaryrepo.Repo{ Name: "redis", Repotype: "remote", Pkgtype: "docker", Remoteurl: "https://registry-1.docker.io"}
		repoServiceMock.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(r, nil)
		anyCtx := mock.Anything
		rc := io.NopCloser(strings.NewReader("Hello, world!"))
		rsp := http.Response{
			Header: make(http.Header),
			StatusCode: http.StatusOK,
			Body:       rc,
		}
		registryClientMock.On("DownloadLayer", anyCtx, "redis", digest).Return(rc, &rsp, nil)
		rsp.Header.Set("Content-Type", "application/octet-stream")
		proxy := NewDockerProxyHandler(repoServiceMock, fileStoreMock, registryClientMock)
		vars := map[string]string{
				"repo-name": "docker-remote",
				"namespace": "redis",
				"digest": d,
		}
		req = mux.SetURLVars(req, vars)
		proxy.DownloadLayer(rec, req)
		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "application/octet-stream", res.Header.Get("Content-Type"))
})

}