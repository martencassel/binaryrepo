package client

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func (r *Registry) Pingable() bool {
	return !strings.HasSuffix(r.URL, "gcr.io")
}

func (r *Registry) Ping(ctx context.Context) error {
	url := r.url("/v2/")
	req, err := http.NewRequest("GET", url, nil)
	log.Println(req)
	if err != nil {
		return err
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
