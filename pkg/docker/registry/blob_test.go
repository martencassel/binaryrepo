package registry

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func TestGetBlob(t *testing.T) {
	t.Run("returns a blob from filestore", func(t *testing.T) {
		// Prepare filestore with blob
		blob, err := ioutil.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631")
		if err != nil {
			t.Fatal(err)
		}
		fs := filestore.NewFileStore("/tmp/filestore")
		fs.Remove(digest.NewDigestFromBytes(digest.SHA256, blob))
		_, err = fs.WriteFile(blob)
		if err != nil {
			log.Fatal(err)
		}
		// Setup registry handler
		index := repo.NewRepoIndex()
		registry := NewDockerRegistry(fs, index)
		// Create a request
		req, _ := http.NewRequest("GET", "", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "docker-local",
			"name":      "redis",
			"digest":    "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631",
		}
		req = mux.SetURLVars(req, vars)
		registry.GetBlob(res, req)
		// Check if response matches wanted
		body := res.Body
		want := "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631"
		got := digest.FromBytes(body.Bytes()).String()
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
		assert.Equal(t, res.Header().Get("Content-Type"), "application/octet-stream")
		assert.Equal(t, res.Header().Get("Content-Length"), fmt.Sprintf("%d", len(blob)))
		assert.Equal(t, res.Header().Get("Docker-Content-Digest"), want)
		assert.Equal(t, http.StatusOK, res.Code)
	})
	t.Run("Check a existing layer", func(t *testing.T) {
		// Prepare filestore with blob
		blob, err := ioutil.ReadFile("./testdata/7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631")
		if err != nil {
			t.Fatal(err)
		}
		fs := filestore.NewFileStore("/tmp/filestore")
		fs.Remove(digest.NewDigestFromBytes(digest.SHA256, blob))
		_, err = fs.WriteFile(blob)
		if err != nil {
			log.Fatal(err)
		}
		// Setup registry handler
		index := repo.NewRepoIndex()
		registry := NewDockerRegistry(fs, index)
		// Create a request
		req, _ := http.NewRequest("GET", "", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "docker-local",
			"name":      "redis",
			"digest":    "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631",
		}
		req = mux.SetURLVars(req, vars)
		registry.HeadBlob(res, req)
		// Check if response matches wanted
		want := "sha256:7614ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631"
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, res.Header().Get("Content-Length"), fmt.Sprintf("%d", len(blob)))
		assert.Equal(t, res.Header().Get("Docker-Content-Digest"), want)
	})
	t.Run("Check for a non-existing layer", func(t *testing.T) {
		fs := filestore.NewFileStore("/tmp/filestore")
		// Setup registry handler
		index := repo.NewRepoIndex()
		registry := NewDockerRegistry(fs, index)
		// Create a request
		req, _ := http.NewRequest("GET", "", nil)
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": "docker-local",
			"name":      "redis",
			"digest":    "sha256:1234ae9453d1d87e740a2056257a6de7135c84037c367e1fffa92ae922784631",
		}
		req = mux.SetURLVars(req, vars)
		registry.HeadBlob(res, req)
		// Check if response matches wanted
		want := http.StatusNotFound
		assert.Equal(t, want, res.Code)
	})
	t.Run("Initiate blob upload", func(t *testing.T) {
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		// Setup registry handler
		registry := NewDockerRegistry(fs, index)
		// Create a request
		req, _ := http.NewRequest("GET", "", nil)
		res := httptest.NewRecorder()
		repoName := "docker-local"
		vars := map[string]string{
			"repo-name": repoName,
			"name":      "redis",
		}
		req = mux.SetURLVars(req, vars)
		registry.InitBlobUpload(res, req)
		assert.Contains(t, res.Header().Get("Location"), fmt.Sprintf("/repo/%s/v2/redis/blobs/uploads/", repoName))
		assert.Equal(t, res.Header().Get("Range"), "bytes=0-0")
		assert.Equal(t, res.Header().Get("Content-Length"), "0")
		assert.True(t, IsValidUUID(res.Header().Get("Docker-Upload-UUID")))
		assert.Equal(t, http.StatusAccepted, res.Code)
		_, err := os.Stat("/tmp/filestore/uploads/" + res.Header().Get("Docker-Upload-UUID"))
		if err != nil {
			t.Fatal(err)
		}
		if errors.Is(err, os.ErrNotExist) {
			t.Fatal(err)
		}
	})
	t.Run("Chunked upload", func(t *testing.T) {
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		// Setup registry handler
		registry := NewDockerRegistry(fs, index)
		c1 := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec tristique ultrices erat, in interdum tortor ullamcorper euismod. Donec vel justo vel nisl mattis posuere eu eget enim. Phasellus porttitor ante vitae feugiat posuere. Nam ut metus quis urna placerat faucibus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Vivamus id arcu at lectus gravida condimentum ac vel nisl. Nulla id velit vel neque fermentum pharetra. Duis fringilla justo vel risus lobortis pulvinar. Integer mollis nisi velit, convallis eleifend metus iaculis sit amet. In neque nibh, vehicula vitae ante ut, bibendum tempor quam. Fusce a tortor ut tellus tempus consectetur. Curabitur odio magna, placerat quis enim a, bibendum dictum elit. Mauris accumsan dolor ut consequat tristique. In bibendum eros libero, vitae ornare ipsum suscipit id. Donec ultricies ornare mauris, nec euismod nisi iaculis at. Vestibulum sit amet felis ac lorem lacinia consectetur eget quis velit. Donec varius blandit leo viverra commodo. Cras maximus nec erat in iaculis. Cras mauris orci, tincidunt vitae pretium quis, tristique vel nulla. Sed et ligula sed augue fermentum bibendum eget sit amet ligula. Curabitur mattis est tellus, vitae ullamcorper elit ornare ut. Morbi fermentum erat metus, eu scelerisque diam euismod eget. Pellentesque laoreet, nunc id blandit ullamcorper, dui orci egestas ante, a luctus mi purus ac ipsum. Pellentesque eros massa, lacinia ac nibh sit amet, blandit rhoncus tellus. Donec mollis dictum faucibus. Vestibulum ullamcorper neque risus, non venenatis nulla ullamcorper eu. Suspendisse lobortis condimentum cursus. Proin ac porta ipsum. Vestibulum metus elit, maximus sed mauris tincidunt, semper mollis velit. Pellentesque nunc odio, lacinia id tristique vel, fringilla quis ipsum. Vestibulum ut tempor elit. Ut dui justo, bibendum ut consectetur vitae, laoreet ut ipsum. Etiam eget elit tristique, fermentum ante vitae, dapibus ipsum. Sed rutrum ex sed dolor interdum aliquam. Donec nisi lacus, lobortis vitae sapien eget, cursus bibendum nulla. Curabitur in felis auctor, dictum felis id, condimentum tortor. Suspendisse potenti. Nunc varius, quam sit amet consequat elementum, diam lacus tempor arcu, vel dictum tortor felis eget est. Nulla nec sapien ut odio euismod ullamcorper nec vitae dolor. Nulla pellentesque urna eu risus dignissim, pellentesque pharetra mauris porta. Pellentesque a tortor ante. Sed odio elit, pellentesque sed mauris sit amet, convallis rhoncus magna. Vestibulum malesuada neque sed libero sagittis fringilla. Nulla quis nibh elit. Curabitur eget malesuada ex. Etiam sit amet nibh justo. Etiam pretium, nulla ultrices posuere commodo, turpis sem ullamcorper nisl, at facilisis est dolor ut nisl. Sed luctus lacus et dui ullamcorper, eget faucibus risus dictum. Suspendisse laoreet condimentum arcu, ac interdum lorem tincidunt mollis. Sed aliquet nec ex sit amet fringilla. Curabitur elit mi, convallis vel condimentum id, tincidunt sit amet quam.")
		//c2 := []byte("Nunc quis ornare ligula. Nam ut aliquam felis. Integer nec orci et erat dignissim ultricies vitae eget lorem.")
		//blob := append(c1[:], c2[:]...)
		// Create a request
		repoName := "docker-local"
		body := bytes.NewBuffer(c1)
		req, _ := http.NewRequest("PATCH", "", body)
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": repoName,
			"name":      "redis",
			"uuid":      "3db61366-7256-11ec-9d50-e86a647ebe1b",
		}
		req.Header.Add("Content-Length", fmt.Sprintf("%d", len(c1)))
		req.Header.Add("Content-Range", fmt.Sprintf("%d-%d", 0, len(c1)))
		req.Header.Add("Content-Type", "application/octet-stream")
		req = mux.SetURLVars(req, vars)
		registry.UploadBlobChunk(res, req)
		assert.Contains(t, res.Header().Get("Location"), fmt.Sprintf("/repo/%s/v2/redis/blobs/uploads/", repoName))
		assert.Equal(t, http.StatusAccepted, res.Code)
		assert.Equal(t, res.Header().Get("Range"), fmt.Sprintf("0-%s", strconv.Itoa(len(c1))))
		assert.Equal(t, res.Header().Get("Content-Length"), "0")
		assert.True(t, IsValidUUID(res.Header().Get("Docker-Upload-UUID")))
	})
	t.Run("Complete chunk", func(t *testing.T) {
		c1 := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec tristique ultrices erat, in interdum tortor ullamcorper euismod. Donec vel justo vel nisl mattis posuere eu eget enim. Phasellus porttitor ante vitae feugiat posuere. Nam ut metus quis urna placerat faucibus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Vivamus id arcu at lectus gravida condimentum ac vel nisl. Nulla id velit vel neque fermentum pharetra. Duis fringilla justo vel risus lobortis pulvinar. Integer mollis nisi velit, convallis eleifend metus iaculis sit amet. In neque nibh, vehicula vitae ante ut, bibendum tempor quam. Fusce a tortor ut tellus tempus consectetur. Curabitur odio magna, placerat quis enim a, bibendum dictum elit. Mauris accumsan dolor ut consequat tristique. In bibendum eros libero, vitae ornare ipsum suscipit id. Donec ultricies ornare mauris, nec euismod nisi iaculis at. Vestibulum sit amet felis ac lorem lacinia consectetur eget quis velit. Donec varius blandit leo viverra commodo. Cras maximus nec erat in iaculis. Cras mauris orci, tincidunt vitae pretium quis, tristique vel nulla. Sed et ligula sed augue fermentum bibendum eget sit amet ligula. Curabitur mattis est tellus, vitae ullamcorper elit ornare ut. Morbi fermentum erat metus, eu scelerisque diam euismod eget. Pellentesque laoreet, nunc id blandit ullamcorper, dui orci egestas ante, a luctus mi purus ac ipsum. Pellentesque eros massa, lacinia ac nibh sit amet, blandit rhoncus tellus. Donec mollis dictum faucibus. Vestibulum ullamcorper neque risus, non venenatis nulla ullamcorper eu. Suspendisse lobortis condimentum cursus. Proin ac porta ipsum. Vestibulum metus elit, maximus sed mauris tincidunt, semper mollis velit. Pellentesque nunc odio, lacinia id tristique vel, fringilla quis ipsum. Vestibulum ut tempor elit. Ut dui justo, bibendum ut consectetur vitae, laoreet ut ipsum. Etiam eget elit tristique, fermentum ante vitae, dapibus ipsum. Sed rutrum ex sed dolor interdum aliquam. Donec nisi lacus, lobortis vitae sapien eget, cursus bibendum nulla. Curabitur in felis auctor, dictum felis id, condimentum tortor. Suspendisse potenti. Nunc varius, quam sit amet consequat elementum, diam lacus tempor arcu, vel dictum tortor felis eget est. Nulla nec sapien ut odio euismod ullamcorper nec vitae dolor. Nulla pellentesque urna eu risus dignissim, pellentesque pharetra mauris porta. Pellentesque a tortor ante. Sed odio elit, pellentesque sed mauris sit amet, convallis rhoncus magna. Vestibulum malesuada neque sed libero sagittis fringilla. Nulla quis nibh elit. Curabitur eget malesuada ex. Etiam sit amet nibh justo. Etiam pretium, nulla ultrices posuere commodo, turpis sem ullamcorper nisl, at facilisis est dolor ut nisl. Sed luctus lacus et dui ullamcorper, eget faucibus risus dictum. Suspendisse laoreet condimentum arcu, ac interdum lorem tincidunt mollis. Sed aliquet nec ex sit amet fringilla. Curabitur elit mi, convallis vel condimentum id, tincidunt sit amet quam.")
		c2 := []byte("Nunc quis ornare ligula. Nam ut aliquam felis. Integer nec orci et erat dignissim ultricies vitae eget lorem.")
		fs := filestore.NewFileStore("/tmp/filestore")
		index := repo.NewRepoIndex()
		// Setup registry handler
		registry := NewDockerRegistry(fs, index)
		log.Print(registry)
		// Create a request
		repoName := "docker-local"
		body := bytes.NewBuffer(c2)
		req, _ := http.NewRequest("PUT", "", body)
		res := httptest.NewRecorder()
		vars := map[string]string{
			"repo-name": repoName,
			"name":      "redis",
		}
		req.Header.Add("Content-Length", fmt.Sprintf("%d", len(c2)))
		req.Header.Add("Content-Range", fmt.Sprintf("%d-%d", len(c1), len(c1)+len(c2)-1))
		req = mux.SetURLVars(req, vars)
		registry.CompleteUpload(res, req)
		assert.Equal(t, http.StatusAccepted, res.Code)
		assert.Contains(t, res.Header().Get("Location"), fmt.Sprintf("/repo/%s/v2/redis/blobs/uploads/", repoName))
		assert.Equal(t, res.Header().Get("Content-Length"), "0")
		assert.True(t, IsValidUUID(res.Header().Get("Docker-Upload-UUID")))
	})
}
