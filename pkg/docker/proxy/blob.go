package dockerproxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	log.Info().Msgf("%s %s", req.Method, req.URL.Path)
	log.Info().Msgf("Repo Name: %s, Namespace: %s, Digest: %s, Namespace 1: %s, Namespace 2: %s", opt.repoName, opt.namespace, opt.digest, opt.namespace1, opt.namespace2)
}

func (p *DockerProxyApp) DownloadLayer(w http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("%s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	if opt.repoName == "" {
		log.Info().Msgf("No repo name")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		log.Info().Msgf("Repo %s was not found", opt.repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	PrintOptions(req, opt)
	log.Info().Msgf("Digest: %s\n", opt.digest)
	// Check if digest exists in filestore, if so
	// then read file and write it to response writer
	digest, err := digest.Parse(opt.digest)
	if err != nil {
		log.Info().Msgf("Digest is invalid %s", err)
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
			log.Info().Msgf("Error writing to response writer %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	log.Print(digest, opt)
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
		log.Error().Msgf("Error getting digest: %s\n", err)
	}

	log.Info().Msgf("DownloadLayer: Status: %s\n", resp.Status)

	if resp.StatusCode == http.StatusTemporaryRedirect {
		log.Info().Msgf("Redirecting to: %s", resp.Header.Get("Location"))
	}
	// Save file to cache

	/*var bodyBytes []byte
	if resp.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	resp.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))*/

	log.Info().Msgf("Content-Length: %s", resp.Header.Get("Content-Length"))

	log.Info().Msg("Here layer response:")
	//	log.Info().Msgf("Layer response body length = %d", len(bodyBytes))

	if resp.Request.Body != nil {
		out, err := os.Create(fmt.Sprintf("/tmp/%s", opt.digest))
		if err != nil {
			log.Error().Msgf("Error creating file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer out.Close()
		n, err := io.Copy(out, resp.Body)
		if err != nil {
			log.Error().Msgf("Error copying file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info().Msgf("Copied %d bytes to file %s \n", n, out.Name())
		bodyBytes, _ := ioutil.ReadFile(out.Name())
		resp.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if err != nil {
			log.Error().Msgf("Error reading response body: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := ioutil.ReadFile(out.Name())
		if err != nil {
			log.Error().Msgf("Error reading file: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = p.fs.WriteFile(b)
		if err != nil {
			log.Error().Msgf("Error writing to cache: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	copyResponse(w, resp)
}
