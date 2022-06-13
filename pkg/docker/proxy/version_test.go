package dockerproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVersionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8081/repo/docker-remote/v2/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	rs := mocks.NewRepoStoreMock()
	fs := mocks.NewFileStoreMock()
	regClient := mocks.NewRegistryClientMock()
	rs.On("Get", mock.Anything, "docker-remote").Return(&binaryrepo.Repo{ Name: "docker-remote", Repotype: "remote", Pkgtype: "docker", Remoteurl: "https://registry-1.docker.io", Anonymous: false })
	p := NewDockerProxyHandler(rs, fs, regClient)
	vars := map[string]string{
		"repo-name": "docker-remote",
	}
	req = mux.SetURLVars(req, vars)
	p.VersionHandler(rec, req)
	res := rec.Result()
	assert.Equal(t, 200, res.StatusCode)
}