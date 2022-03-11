package dockerproxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	regclient "github.com/martencassel/binaryrepo/pkg/docker/client"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
	log "github.com/rs/zerolog/log"
)

// PathGetManifest URL.
const PathGetManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathGetManifest2 = "/repo/{repo-name}/v2/{namespace}/{namespace2}/manifests/{reference}"

/*
	Get a manifest
	GET /v2/<name>/manifests/<reference>
*/
func (p *DockerProxyApp) GetManifestHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.getmanifest %s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	_repo := p.index.FindRepo(opt.repoName)
	if _repo == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Username == "" || _repo.Password == "" {
		log.Error().Msgf("Repo %s is not authorized. Add username/password in cmd/binary-repo/run.go", opt.repoName)
		rw.WriteHeader(http.StatusUnauthorized)
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
	if err != nil {
		log.Error().Msgf("Error creating registry client: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.Digest(ctx, regclient.Image{Domain: "docker.io", Path: path, Tag: opt.reference})
	if err != nil {
		log.Error().Msgf("Error getting digest: %s", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	copyResponse(rw, resp)
}

func copyResponse(rw http.ResponseWriter, resp *http.Response) {
	copyHeader(rw.Header(), resp.Header)
	rw.WriteHeader(resp.StatusCode)
	_, err := io.Copy(rw, resp.Body)
	if err != nil {
		log.Error().Msgf("Error copying response: %s", err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// PathHeadManifest URL.
const PathHeadManifest1 = "/repo/{repo-name}/v2/{namespace}/manifests/{reference}"
const PathHeadManifest2 = "/repo/{repo-name}/v2/{namespace1}/{namespace2}/manifests/{reference}"

/*
	Check for a manifest
	HEAD /v2/<name>/manifests/<digest>
*/
func (p *DockerProxyApp) HasManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("proxy.head_manifest %s %s", req.Method, req.URL.Path)
	opt := GetOptions(req)
	var _repo *repo.Repo
	vars := mux.Vars(req)
	repoName := vars["repo-name"]
	_repo = p.index.FindRepo(repoName)
	if _repo == nil {
		log.Error().Msgf("Repo %s was not found", repoName)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	if _repo.Username == "" || _repo.Password == "" {
		log.Error().Msgf("Repo %s is not authorized", opt.repoName)
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
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
		log.Error().Msgf("Error creating registry client: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	path := fmt.Sprintf("library/%s", opt.namespace)
	_, resp, err := r.Digest(ctx, regclient.Image{Domain: "docker.io", Path: path, Tag: opt.reference})
	if err != nil {
		log.Error().Msgf("Error getting digest: %s", err)
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	copyResponse(rw, resp)
}
