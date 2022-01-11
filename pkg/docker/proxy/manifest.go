package dockerproxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
)

// PathGetManifest URL.
const PathGetManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathGetManifest2 = "/repo/{repo-name}/v2/{namespace}/{namespace2}/manifests/{reference}"

// GetManifestHandler implements GET baseURL/repo/v2/namespace/manifests/reference
func (p *DockerProxyApp) GetManifestHandler(w http.ResponseWriter, req *http.Request) {
	opt := GetOptions(req)
	log.Printf("%s %s\n", req.Method, req.URL.Path)
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Proxy the request to the remote server
	manifestDirector := manifestProxyDirector(opt, opt.reference, "registry-1.docker.io")
	manifestModifyResponse := manifestModifyResponse(opt)
	_url, err := url.Parse(_repo.URL)
	if err != nil {
		log.Printf("Url was invalid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errorHandler := func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Ooops error")
		log.Print(err)
	}
	scope := getScope(opt)
	transporter := NewReAuthTransport(_repo, scope, _repo.AccessToken)
	proxy := httputil.NewSingleHostReverseProxy(_url)
	proxy.Director = manifestDirector
	proxy.ModifyResponse = manifestModifyResponse
	proxy.Transport = transporter
	proxy.ErrorHandler = errorHandler
	proxy.ServeHTTP(w, req)
	_repo.AccessToken = transporter.accessToken
}

// PathHeadManifest URL.
const PathHeadManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathHeadManifest2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/manifests/{reference}"

// HeadManifestHandler implements GET baseURL/repo/v2/{namespace}/manifests/{reference}
func (p *DockerProxyApp) HeadManifestHandler(w http.ResponseWriter, req *http.Request) {
	opt := GetOptions(req)
	var _repo *repo.Repo
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo = p.index.FindRepo(repoName)
	if _repo == nil {
		log.Printf("Repo %s was not found", repoName)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	reference := vars["reference"]
	manifestDirector := manifestProxyDirector(opt, reference, "registry-1.docker.io")
	manifestModifyResponse := manifestModifyResponse(opt)
	_url, err := url.Parse(_repo.URL)
	if err != nil {
		log.Printf("Url was invalid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	errorHandler := func(http.ResponseWriter, *http.Request, error) {
		log.Printf("Ooops error")
	}
	scope := getScope(opt)
	transporter := NewReAuthTransport(_repo, scope, _repo.AccessToken)
	proxy := httputil.NewSingleHostReverseProxy(_url)
	proxy.Director = manifestDirector
	proxy.ModifyResponse = manifestModifyResponse
	proxy.Transport = transporter
	proxy.ErrorHandler = errorHandler
	proxy.ServeHTTP(w, req)
	_repo.AccessToken = transporter.accessToken
}

func manifestProxyDirector(opt HandlerOptions, reference, upstreamHost string) func(*http.Request) {
	f := func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = upstreamHost
		req.Host = req.URL.Host
		path := req.URL.Path
		components := strings.Split(path, "/")
		imageName := components[4]
		newImageName := fmt.Sprintf("library/%s", imageName)
		newUrlPath := strings.Replace(path, imageName, newImageName, 1)
		newUrlPath = strings.Replace(newUrlPath, fmt.Sprintf("/%s/%s", components[1], components[2]), "", 1)
		//	log.Printf("Rewriting URL %s to %s", req.URL.Path, newUrlPath)
		req.URL.Path = newUrlPath
	}
	return f
}

func manifestModifyResponse(opt HandlerOptions) func(resp *http.Response) error {
	f := func(resp *http.Response) error {
		if resp.StatusCode == http.StatusUnauthorized {
			log.Print("Got Unauthorized!!")
		}
		var b []byte
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}
		body := ioutil.NopCloser(bytes.NewReader(b))
		resp.Body = body
		resp.ContentLength = int64(len(b))
		resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
		return nil
	}
	return f
}
