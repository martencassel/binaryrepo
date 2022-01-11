package dockerproxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestGetBlobNotInCache(t *testing.T) {
	log.Printf("TestGetBlob")
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
	p := NewDockerProxyApp()
	p.index.AddRepo(repo.Repo{
		ID:      1,
		Name:    "docker-remote",
		Type:    repo.Remote,
		PkgType: repo.Docker,
		URL:     "https://registry-1.docker.io",
	})
	err = p.fs.Remove(digest)
	if err != nil {
		t.Fatal(err)
	}
	vars := map[string]string{
		"repo-name": "docker-remote",
		"namespace": "redis",
		"digest":    d,
	}
	req = mux.SetURLVars(req, vars)
	p.GetBlobHandler(rec, req)
	res := rec.Result()
	assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
	assert.Equal(t, "/repo/docker-remote/v2/blob/sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631", res.Header.Get("Location"))
}

func TestGetBlobInCache(t *testing.T) {
	err := os.MkdirAll("/tmp/filestore/13/", 0777)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile("./testdata/1e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("/tmp/filestore/13/1e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f", b, 0777)
	if err != nil {
		t.Fatal(err)
	}
	d := "sha256:131e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f"
	req, err := http.NewRequest("GET", "http://localhost:8081/repo/docker-remote/v2/redis/blobs/sha256:131e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	p := NewDockerProxyApp()
	vars := map[string]string{
		"repo-name": "docker-remote",
		"namespace": "redis",
		"digest":    d,
	}
	req = mux.SetURLVars(req, vars)
	p.GetBlobHandler(rec, req)
	res := rec.Result()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "201", res.Header.Get("Content-Length"))
	assert.Equal(t, "sha256:131e916e1a28efae8398ab187eaf75683c6c7ebc71e90f780e19a95465dfd52f", digest.FromBytes(b).String())
}
