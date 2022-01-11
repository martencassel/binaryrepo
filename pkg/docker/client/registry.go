package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema2"
)

type Registry struct {
	URL      string
	Domain   string
	Username string
	Password string
	Client   *http.Client
	Opt      Opt
}

var reProtocol = regexp.MustCompile("^https?://")

type Opt struct {
	Domain   string
	SkipPing bool
	Timeout  time.Duration
	NonSSL   bool
	Headers  map[string]string
}

func New(ctx context.Context, auth AuthConfig, opt Opt) (*Registry, error) {
	transport := http.DefaultTransport
	return newFromTransport(ctx, auth, transport, opt)
}

func newFromTransport(ctx context.Context, auth AuthConfig, transport http.RoundTripper, opt Opt) (*Registry, error) {
	if len(opt.Domain) < 1 || opt.Domain == "docker.io" {
		opt.Domain = auth.ServerAddress
	}
	url := strings.TrimSuffix(opt.Domain, "/")
	authURL := strings.TrimSuffix(auth.ServerAddress, "/")

	if !reProtocol.MatchString(url) {
		if !opt.NonSSL {
			url = "https://" + url
		} else {
			url = "http://" + url
		}
	}
	tokenTransport := &TokenTransport{
		Transport: transport,
		Username:  auth.Username,
		Password:  auth.Password,
		Scope:     auth.Scope,
	}
	basicAuthTransport := &BasicTransport{
		Transport: tokenTransport,
		URL:       authURL,
		Username:  auth.Username,
		Password:  auth.Password,
	}
	errorTransport := &ErrorTransport{
		Transport: basicAuthTransport,
	}
	customTransport := &CustomTransport{
		Transport: errorTransport,
		Headers:   opt.Headers,
	}
	registry := &Registry{
		URL:    url,
		Domain: reProtocol.ReplaceAllString(url, ""),
		Client: &http.Client{
			Timeout:   opt.Timeout,
			Transport: customTransport,
		},
		Username: auth.Username,
		Password: auth.Password,
	}
	if registry.Pingable() && !opt.SkipPing {
		if err := registry.Ping(ctx); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

func (r *Registry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *Registry) getJSON(ctx context.Context, url string, response interface{}) (http.Header, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	switch response.(type) {
	case *schema2.Manifest:
		req.Header.Add("Accept", schema2.MediaTypeManifest)
	case *manifestlist.ManifestList:
		req.Header.Add("Accept", manifestlist.MediaTypeManifestList)
	}
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("registry.registry resp.Status=%s", resp.Status)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return resp.Header, nil
}
