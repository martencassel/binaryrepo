package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/opencontainers/go-digest"
	log "github.com/rs/zerolog/log"
)

func (r *Registry) Digest(ctx context.Context, image Image) (digest.Digest, *http.Response, error) {
	if len(image.Digest) > 1 {
		return image.Digest, nil, nil
	}
	url := r.url("/v2/%s/manifests/%s", image.Path, image.Tag)
	log.Info().Msgf("manifests.get url=%s repository=%s ref=%s",
		url, image.Path, image.Tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Add("Accept", schema2.MediaTypeManifest)
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err != nil {
		/*		if resp != nil {
				defer resp.Body.Close()
			}*/
		return "", resp, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return "", resp, fmt.Errorf("got status code: %d", resp.StatusCode)
	}
	d, err := digest.Parse(resp.Header.Get("Docker-Content-Digest"))
	if err != nil {
		return "", resp, err
	}
	return d, resp, nil
}
