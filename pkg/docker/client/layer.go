package client

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/opencontainers/go-digest"
	log "github.com/rs/zerolog/log"
)

// DownloadLayer downloads a specific layer by digest for a repository.
func (r *Registry) DownloadLayer(ctx context.Context, repository string, digest digest.Digest) (io.ReadCloser, *http.Response, error) {
	url := r.url("/v2/%s/blobs/%s", repository, digest)
	log.Info().Msgf("registry.layer.download url=%s repository=%s digest=%s", url, repository, digest)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, resp, err
	}
	return resp.Body, resp, nil
}

// HasLayer returns if the registry contains the specific digest for a repository.
func (r *Registry) HasLayer(ctx context.Context, repository string, digest digest.Digest) (bool, *http.Response, error) {
	checkURL := r.url("/v2/%s/blobs/%s", repository, digest)
	log.Info().Msgf("registry.layer.check url=%s repository=%s digest=%s", checkURL, repository, digest)
	req, err := http.NewRequest("HEAD", checkURL, nil)
	if err != nil {
		return false, nil, err
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err == nil {
		//defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK, resp, nil
	}
	urlErr, ok := err.(*url.Error)
	if !ok {
		return false, resp, err
	}
	httpErr, ok := urlErr.Err.(*httpStatusError)
	if !ok {
		return false, resp, err
	}
	if httpErr.Response.StatusCode == http.StatusNotFound {
		return false, resp, nil
	}
	return false, resp, err
}
