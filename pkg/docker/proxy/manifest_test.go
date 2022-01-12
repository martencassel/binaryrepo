package dockerproxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// GET	 /v2/<name>/manifests/<reference>
func TestGetManifest(t *testing.T) {

	t.Run("Fetch a manifest", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8081/repo/docker-remote/v2/redis/manifests/latest", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		p := NewDockerProxyApp()
		vars := map[string]string{
			"repo-name":  "docker-remote",
			"namespace":  "redis",
			"namespace1": "",
			"namespace2": "",
			"reference":  "latest",
		}
		req = mux.SetURLVars(req, vars)
		p.GetManifestHandler(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Status code is not OK: %d", res.StatusCode)
		}
		assert.Equal(t, "13646", res.Header.Get("Content-Length"))
		assert.Equal(t, "application/vnd.docker.distribution.manifest.v1+prettyjws", res.Header.Get("Content-Type"))
		assert.Equal(t, "sha256:aa424bde415a742dead1a210b9b2481753307a97016a0e03840fd6c2e26bc75b", res.Header.Get("docker-content-digest"))
		assert.Equal(t, "registry/2.0", res.Header.Get("Docker-Distribution-Api-Version"))
		assert.Equal(t, "\"sha256:aa424bde415a742dead1a210b9b2481753307a97016a0e03840fd6c2e26bc75b\"", res.Header.Get("Etag"))
	})

}

// HEAD /v2/<name>/manifests/<reference>
func TestManifestExists(t *testing.T) {

	t.Run("Test 2", func(t *testing.T) {
		req, err := http.NewRequest("HEAD", "http://localhost:8081/repo/docker-remote/v2/redis/manifests/latest", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		p := NewDockerProxyApp()
		vars := map[string]string{
			"repo-name":  "docker-remote",
			"namespace":  "redis",
			"namespace1": "",
			"namespace2": "",
			"reference":  "latest",
		}
		req = mux.SetURLVars(req, vars)
		p.HeadManifestHandler(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Status code is not OK: %d", res.StatusCode)
		}
		assert.Equal(t, "application/vnd.docker.distribution.manifest.v1+prettyjws", res.Header.Get("Content-Type"))
		assert.Equal(t, "sha256:aa424bde415a742dead1a210b9b2481753307a97016a0e03840fd6c2e26bc75b", res.Header.Get("docker-content-digest"))
		assert.Equal(t, "registry/2.0", res.Header.Get("Docker-Distribution-Api-Version"))
		assert.Equal(t, "\"sha256:aa424bde415a742dead1a210b9b2481753307a97016a0e03840fd6c2e26bc75b\"", res.Header.Get("Etag"))
	})

	t.Run("Test 3", func(t *testing.T) {
		req, err := http.NewRequest("HEAD", "http://localhost:8081/repo/docker-remote/v2/redis/manifests/latest", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		p := NewDockerProxyApp()
		vars := map[string]string{
			"repo-name":  "not-found",
			"namespace":  "redis",
			"namespace1": "",
			"namespace2": "",
			"reference":  "latest",
		}
		req = mux.SetURLVars(req, vars)
		p.GetManifestHandler(rec, req)
		res := rec.Result()
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
