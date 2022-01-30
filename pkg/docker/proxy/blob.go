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
func (p *DockerProxyApp) DownloadLayer(w http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.downloadlayer %s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	if opt.repoName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(opt.digest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	blobExists := p.fs.Exists(digest)
	if blobExists {
		b, err := p.fs.ReadFile(digest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		_, err = w.Write(b)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	ctx := context.Background()
	scope := fmt.Sprintf("repository:library/%s:pull", opt.namespace)
	r, err := p.NewRegistryClient("docker.io", _repo.Username, _repo.Password, scope, _repo.URL)
	if err != nil {
		log.Error().Msgf("Error creating registry client: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.DownloadLayer(ctx, path, digest)
	if err != nil {
		log.Error().Msgf("Error getting digest: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Msgf("Error reading response body: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, err := p.fs.WriteFile(bodyBytes)
		if err != nil {
			log.Error().Msgf("Error writing to cache: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Restore the io.ReadCloser to its original state
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	copyResponse(w, resp)
}

const PathHeadBlob1 = "/repo/{repo-name}/v2/{namespace}/blobs/{digest}"
const PathHeadBlob2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/blobs/{digest}"

func (p *DockerProxyApp) LayerPut(w http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.LayerPut %s %s", req.Method, req.URL.Path)
	data := "BLOB_UNKNOWN"
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error().Msgf("Error encoding JSON: %s", err)
	}
	w.WriteHeader(http.StatusOK)
}

/*
	Check for layer
	HEAD /v2/<name>/blobs/<digest>
*/
func (p *DockerProxyApp) HasLayer(w http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.hasLayer %s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	if opt.repoName == "" {
		log.Error().Msg("No repo name")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		log.Error().Msgf("Repo %s was not found", opt.repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	digest, err := digest.Parse(opt.digest)
	if err != nil {
		log.Info().Msgf("Digest is invalid %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	blobExists := p.fs.Exists(digest)
	if blobExists {
		w.Header().Set("Location", fmt.Sprintf("/v2/%s/blobs/%s", opt.namespace, opt.digest))
		w.WriteHeader(http.StatusOK)
		return
	}
	ctx := context.Background()
	scope := fmt.Sprintf("repository:library/%s:pull", opt.namespace)
	r, err := p.NewRegistryClient("docker.io", _repo.Username, _repo.Password, scope, _repo.URL)
	if err != nil {
		log.Error().Msgf("Error creating registry client: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.HasLayer(ctx, path, digest)
	if err != nil {
		log.Error().Msgf("Error getting digest: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Info().Msgf("%v", resp)
	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
