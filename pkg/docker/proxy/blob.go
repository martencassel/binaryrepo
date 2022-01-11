package dockerproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
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

func (p *DockerProxyApp) GetBlobHandler(w http.ResponseWriter, req *http.Request) {
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
	manifestDirector := manifestProxyDirector(opt, "", "registry-1.docker.io")
	blobModifyResponse := getBlobModifyResponse(p.fs, opt)
	_url, err := url.Parse(_repo.URL)
	if err != nil {
		log.Printf("Url was invalid")
		w.WriteHeader(http.StatusInternalServerError)
	}
	errorHandler := func(rw http.ResponseWriter, r *http.Request, err error) {
		log.Print(err)
	}
	/*
		Set scope for auth request
	*/
	scope := getScope(opt)
	transporter := NewReAuthTransport(_repo, scope, _repo.AccessToken)
	proxy := httputil.NewSingleHostReverseProxy(_url)
	proxy.Director = manifestDirector
	proxy.ModifyResponse = blobModifyResponse
	proxy.Transport = transporter
	proxy.ErrorHandler = errorHandler
	proxy.ServeHTTP(w, req)
	//	_repo.AccessToken = transporter.accessToken
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

func FetchBlob(urlIn string) (b []byte, err error) {
	log.Printf("FetchBlob\n")
	log.Printf("%s %s", http.MethodGet, urlIn)
	req, err := http.NewRequest("GET", urlIn, nil)
	if err != nil {
		return nil, err
	}
	_url, _ := url.Parse(urlIn)
	req.Header.Set("host", _url.Host)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Connection", "close")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Status == "401" {
		log.Printf("401 Unauthorized")
		return nil, fmt.Errorf("Unauthorized")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	return body, nil
}

func getBlobModifyResponse(fs *filestore.FileStore, opt HandlerOptions) func(resp *http.Response) error {
	f := func(resp *http.Response) (err error) {
		log.Printf("%s %s %s", resp.Request.Method, resp.Request.URL.Path, resp.Status)
		log.Printf("getBlobModifyResponse(%s)", resp.Request.URL.Path)
		PrintOptions(resp.Request, opt)
		log.Printf("Location: %s", resp.Header.Get("Location"))
		digest := digest.Digest(opt.digest)
		exists := fs.Exists(digest)
		log.Print("Exists: %w", exists)
		// If the blobs exists, set Location and redirect to ServeBlobHandler
		if exists {
			newLocation := fmt.Sprintf("/repo/%s/v2/blob/%s", opt.repoName, digest)
			log.Print("Redirecting to ", newLocation)
			resp.Header.Set("Location", newLocation)
			resp.StatusCode = http.StatusTemporaryRedirect
			return
		}
		if resp.StatusCode == http.StatusTemporaryRedirect {
			b, err := FetchBlob(resp.Header.Get("location"))
			if err != nil {
				log.Printf("Err: %s", err)
				return err
			}
			digest, err := fs.WriteFile(b)
			if err != nil {
				log.Printf("Err: %s", err)
				return err
			}
			log.Printf("Stored blob %s", digest)
			newLocation := fmt.Sprintf("/repo/%s/v2/blob/%s", opt.repoName, digest)
			log.Print("Redirecting to ", newLocation)
			resp.Header.Set("Location", newLocation)
			resp.StatusCode = http.StatusTemporaryRedirect
		}
		return
	}
	return f
}
