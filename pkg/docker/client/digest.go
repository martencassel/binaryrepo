package client

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/opencontainers/go-digest"
)

func (r *Registry) Digest(ctx context.Context, image Image) (digest.Digest, error) {
	if len(image.Digest) > 1 {
		return image.Digest, nil
	}
	url := r.url("/v2/%s/manifests/%s", image.Path, image.Tag)
	log.Printf("manifests.get url=%s repository=%s ref=%s",
		url, image.Path, image.Tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", schema2.MediaTypeManifest)
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err != nil {
		defer resp.Body.Close()
		return "", err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return "", fmt.Errorf("got status code: %d", resp.StatusCode)
	}
	return digest.Parse(resp.Header.Get("Docker-Content-Digest"))
}
