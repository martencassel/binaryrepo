package dockerproxy

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/docker/client"

	"github.com/opencontainers/go-digest"
	"github.com/rs/zerolog/log"
)

type HandlerOptions struct {
	repoName   string
	digest     string
	namespace  string
	namespace1 string
	namespace2 string
	reference  string
}

func GetOptions(req *http.Request) HandlerOptions {
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	in_digest := vars["digest"]
	namespace := vars["namespace"]
	namespace1 := vars["namespace1"]
	namespace2 := vars["namespace2"]
	reference := vars["reference"]
	return HandlerOptions{
			repoName:   repoName,
			digest:     in_digest,
			namespace:  namespace,
			namespace1: namespace1,
			namespace2: namespace2,
			reference:  reference,
	}
}
/*
	Get blob
	GET /v2/<name>/blobs/<digest>
*/
func (proxy * DockerProxyHandler) DownloadLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.downloadlayer %s %s", req.Method, req.URL.Path)
	// Check if the repo exists
	opt := GetOptions(req)
	repo, err := proxy.repoStore.Get(context.Background(), opt.repoName)
	if err != nil {
		log.Fatal().Msgf("proxy.downloadlayer %s %s", req.Method, req.URL.Path)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	// Parse the digest of the blob
	digest, err := digest.Parse(opt.digest)
	if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	// If the blob is found in the cache, return it
	blobExists := proxy.fs.Exists(digest)
	if blobExists {
		b, err := proxy.fs.ReadFile(digest)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Length", strconv.Itoa(len(b)))
		rw.Header().Set("Content-Type", "application/octet-stream")
		_, err = rw.Write(b)
		if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// If the blob is not found in the cache, download it from upstream.
	scope := fmt.Sprintf("repository:library/%s:pull", opt.namespace)
	authConfig := &client.AuthConfig {
		Username:  repo.RemoteUsername,
		Password:  repo.RemotePassword,
		Scope:     scope,
	}
	proxy.registryClient.SetConfig(repo.Remoteurl, "docker.io", authConfig)
	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.WriteHeader(http.StatusOK)
}