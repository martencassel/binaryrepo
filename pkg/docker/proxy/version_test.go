package dockerproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestVersionHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8081/repo/docker-remote/v2/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	p := NewDockerProxyApp()
	vars := map[string]string{
		"repo-name": "docker-remote",
	}
	req = mux.SetURLVars(req, vars)
	p.VersionHandler(rec, req)
	res := rec.Result()
	assert.Equal(t, 200, res.StatusCode)
}
