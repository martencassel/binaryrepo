package client

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/opencontainers/go-digest"
)

// DownloadLayer downloads a specific layer by digest for a repository.
func (r *Registry) DownloadLayer(ctx context.Context, repository string, digest digest.Digest) (io.ReadCloser, error) {
	url := r.url("/v2/%s/blobs/%s", repository, digest)
	log.Printf("registry.layer.download url=%s repository=%s digest=%s", url, repository, digest)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// HasLayer returns if the registry contains the specific digest for a repository.
func (r *Registry) HasLayer(ctx context.Context, repository string, digest digest.Digest) (bool, error) {
	checkURL := r.url("/v2/%s/blobs/%s", repository, digest)
	log.Printf("registry.layer.check url=%s repository=%s digest=%s", checkURL, repository, digest)

	req, err := http.NewRequest("HEAD", checkURL, nil)
	if err != nil {
		return false, err
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err == nil {
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK, nil
	}

	urlErr, ok := err.(*url.Error)
	if !ok {
		return false, err
	}
	httpErr, ok := urlErr.Err.(*httpStatusError)
	if !ok {
		return false, err
	}
	if httpErr.Response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, err
}
