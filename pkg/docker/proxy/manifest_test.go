package dockerproxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/stretchr/testify/assert"
)

// GET	 /v2/<name>/manifests/<reference>
func TestGetManifest(t *testing.T) {

	t.Run("Fetch a manifest", func(t *testing.T) {
		b, err := os.ReadFile("./testdata/1e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile("/tmp/filestore/13/1e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f", b, 0777)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("GET", "http://localhost:8081/repo/docker-remote/v2/redis/manifests/latest", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		p := NewDockerProxyApp()
		hubUser := os.Getenv("DOCKERHUB_USERNAME")
		hubPass := os.Getenv("DOCKERHUB_PASSWORD")
		p.index.AddRepo(repo.Repo{
			ID:       1,
			Name:     "docker-remote",
			Type:     repo.Remote,
			PkgType:  repo.Docker,
			URL:      "https://registry-1.docker.io",
			Username: hubUser,
			Password: hubPass,
		})
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
		assert.Equal(t, "1573", res.Header.Get("Content-Length"))
		assert.Equal(t, "application/vnd.docker.distribution.manifest.v2+json", res.Header.Get("Content-Type"))
		assert.Equal(t, "sha256:563888f63149e3959860264a1202ef9a644f44ed6c24d5c7392f9e2262bd3553", res.Header.Get("docker-content-digest"))
		assert.Equal(t, "registry/2.0", res.Header.Get("Docker-Distribution-Api-Version"))
		assert.Equal(t, "\"sha256:563888f63149e3959860264a1202ef9a644f44ed6c24d5c7392f9e2262bd3553\"", res.Header.Get("Etag"))
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
		hubUser := os.Getenv("DOCKERHUB_USERNAME")
		hubPass := os.Getenv("DOCKERHUB_PASSWORD")
		p.index.AddRepo(repo.Repo{
			ID:       1,
			Name:     "docker-remote",
			Type:     repo.Remote,
			PkgType:  repo.Docker,
			URL:      "https://registry-1.docker.io",
			Username: hubUser,
			Password: hubPass,
		})
		vars := map[string]string{
			"repo-name":  "docker-remote",
			"namespace":  "redis",
			"namespace1": "",
			"namespace2": "",
			"reference":  "latest",
		}
		req = mux.SetURLVars(req, vars)
		p.HasManifest(rec, req)
		res := rec.Result()
		if res.StatusCode != http.StatusOK {
			t.Errorf("Status code is not OK: %d", res.StatusCode)
		}
		assert.Equal(t, "application/vnd.docker.distribution.manifest.v2+json", res.Header.Get("Content-Type"))
		assert.Equal(t, "sha256:563888f63149e3959860264a1202ef9a644f44ed6c24d5c7392f9e2262bd3553", res.Header.Get("docker-content-digest"))
		assert.Equal(t, "registry/2.0", res.Header.Get("Docker-Distribution-Api-Version"))
		assert.Equal(t, "\"sha256:563888f63149e3959860264a1202ef9a644f44ed6c24d5c7392f9e2262bd3553\"", res.Header.Get("Etag"))
	})

	t.Run("Test 3", func(t *testing.T) {
		req, err := http.NewRequest("HEAD", "http://localhost:8081/repo/docker-remote/v2/redis/manifests/latest", nil)
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		p := NewDockerProxyApp()
		hubUser := os.Getenv("DOCKERHUB_USERNAME")
		hubPass := os.Getenv("DOCKERHUB_PASSWORD")
		p.index.AddRepo(repo.Repo{
			ID:       1,
			Name:     "docker-remote",
			Type:     repo.Remote,
			PkgType:  repo.Docker,
			URL:      "https://registry-1.docker.io",
			Username: hubUser,
			Password: hubPass,
		})
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
