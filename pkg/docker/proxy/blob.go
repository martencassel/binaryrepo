package dockerproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/opencontainers/go-digest"
	log "github.com/rs/zerolog/log"
)

// PathGetBlob URL.
const PathGetBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/{digest}"
const PathGetBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/{digest}"

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

func PrintOptions(req *http.Request, opt HandlerOptions) {
}

/*
	Get blob
	GET /v2/<name>/blobs/<digest>
*/
func (p *DockerProxyApp) DownloadLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.downloadlayer %s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	if opt.repoName == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(opt.digest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	blobExists := p.fs.Exists(digest)
	if blobExists {
		b, err := p.fs.ReadFile(digest)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		rw.Header().Set("Content-Length", strconv.Itoa(len(b)))
		_, err = rw.Write(b)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	ctx := context.Background()
	scope := fmt.Sprintf("repository:library/%s:pull", opt.namespace)
	r, err := p.NewRegistryClient("docker.io", _repo.Username, _repo.Password, scope, _repo.URL)
	if err != nil {
		log.Error().Msgf("Error creating registry client: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.DownloadLayer(ctx, path, digest)
	if err != nil {
		log.Error().Msgf("Error getting digest: %s", err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("Error reading response body: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
		_, err := p.fs.WriteFile(bodyBytes)
		if err != nil {
			log.Error().Msgf("Error writing to cache: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	copyResponse(rw, resp)
}

const PathHeadBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/{digest}"
const PathHeadBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/{digest}"

func (p *DockerProxyApp) LayerPut(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.LayerPut %s %s", req.Method, req.URL.Path)
	data := "BLOB_UNKNOWN"
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		log.Error().Msgf("Error encoding JSON: %s", err)
	}
	rw.WriteHeader(http.StatusOK)
}

/*
	Check for layer
	HEAD /v2/<name>/blobs/<digest>
*/
func (p *DockerProxyApp) HasLayer(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.hasLayer %s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	if opt.repoName == "" {
		log.Error().Msg("No repo name")
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		log.Error().Msgf("Repo %s was not found", opt.repoName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(opt.digest)
	if err != nil {
		log.Info().Msgf("Digest is invalid %s", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	blobExists := p.fs.Exists(digest)
	if blobExists {
		rw.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/%s", opt.namespace, opt.digest))
		rw.WriteHeader(http.StatusOK)
		return
	}
	ctx := context.Background()
	scope := fmt.Sprintf("repository:library/%s:pull", opt.namespace)
	r, err := p.NewRegistryClient("docker.io", _repo.Username, _repo.Password, scope, _repo.URL)
	if err != nil {
		log.Error().Msgf("Error creating registry client: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.HasLayer(ctx, path, digest)
	if err != nil {
		log.Error().Msgf("Error getting digest: %s", err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	log.Info().Msgf("%v", resp)
	if resp.StatusCode == http.StatusOK {
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}
