package client

import (
	"context"
	"net/http"
	"strings"
)

func (r *registryClient) Pingable() bool {
	return !strings.HasSuffix(r.URL, "gcr.io")
}

func (r *registryClient) Ping(ctx context.Context) error {
	url := r.url("/v2/")
	req, err := http.NewRequest("GET", url, nil)
	////log.Info().Msgf("registry.ping url=%s", url)
	if err != nil {
		return err
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
