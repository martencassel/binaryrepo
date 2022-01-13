package dockerproxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	regclient "github.com/martencassel/binaryrepo/pkg/docker/client"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/opencontainers/go-digest"
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
	log.Printf("%s %s", req.Method, req.URL.Path)
	log.Printf("Repo Name: %s, Namespace: %s, Digest: %s, Namespace 1: %s, Namespace 2: %s", opt.repoName, opt.namespace, opt.digest, opt.namespace1, opt.namespace2)
}

func (p *DockerProxyApp) DownloadLayer(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	if opt.repoName == "" {
		log.Printf("No repo name")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var _repo *repo.Repo
	if p.index.FindRepo(opt.repoName) == nil {
		log.Printf("Repo %s was not found", opt.repoName)
		w.WriteHeader(http.StatusNotFound)
	}
	PrintOptions(req, opt)
	log.Printf("Digest: %s\n", opt.digest)
	// Check if digest exists in filestore, if so
	// then read file and write it to response writer
	digest, err := digest.Parse(opt.digest)
	if err != nil {
		log.Printf("Digest is invalid %s", err)
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
			log.Printf("Error writing to response writer %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// If digest does not exist in filestore, then
	// Proxy the request to the remote server
	// and write the response to response writer.
	_repo = p.index.FindRepo(opt.repoName)
	if _repo == nil {
		log.Printf("Repo %s was not found", opt.repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Print(digest, opt)
	ctx := context.Background()
	scope := fmt.Sprintf("repository:library/%s:pull", opt.namespace)
	r, err := regclient.New(ctx, regclient.AuthConfig{
		Username:      _repo.Username,
		Password:      _repo.Password,
		Scope:         scope,
		ServerAddress: _repo.URL,
	}, regclient.Opt{
		Domain:   "docker.io",
		SkipPing: false,
		Timeout:  time.Second * 120,
		NonSSL:   false,
		Insecure: false,
	})
	if err != nil {
		log.Printf("Error creating registry client: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.DownloadLayer(ctx, path, digest)
	if err != nil {
		log.Printf("Error getting digest: %s\n", err)
	}
	copyResponse(w, resp)
}

// PathGetBlob URL.
const PathServeBlobURL = "/repo/{repo-name}/v2/blob/{digest}"

// GET /repo/{repo-name}/v2/blob/{digest}
func (p *DockerProxyApp) ServeBlobHandler(w http.ResponseWriter, req *http.Request) {
	//	log.Printf("%s %s %s", req.Method, req.URL.Path, req.Response.Status)
	vars := mux.Vars(req)
	in_digest := vars["digest"]
	digest, err := digest.Parse(in_digest)
	log.Printf("%s %s", req.Method, req.URL.Path)
	log.Printf("ServeBlobHandler: %s", in_digest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	blobExists := p.fs.Exists(digest)
	if blobExists {
		log.Printf("Served blob %s", in_digest)
		b, err := p.fs.ReadFile(digest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("Error writing to response writer %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
