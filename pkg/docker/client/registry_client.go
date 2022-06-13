package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	digest "github.com/opencontainers/go-digest"
)

type RegistryClient interface {
	SetConfig(url string, domain string, config *AuthConfig)
	DownloadLayer(ctx context.Context, repository string, digest digest.Digest) (io.ReadCloser, *http.Response, error)
	HasLayer(ctx context.Context, repository string, digest digest.Digest) (bool, *http.Response, error)
	Pingable() bool
	Ping(ctx context.Context) bool
	Manifest(ctx context.Context, repository string, ref string) (distribution.Manifest,error)
	ManifestList(ctx context.Context, repository, ref string) (manifestlist.ManifestList, error)
	ManifestV2(ctx context.Context, repository, ref string) (schema2.Manifest, error)
	ManifestV1(ctx context.Context, repository, ref string) (schema1.SignedManifest, error)
	Token(ctx context.Context, url string) (string, error)
}

var _RegistryClient = &registryClient{}

type registryClient struct {
	URL      string
	Domain   string
	Username string
	Password string
	Scope    string
	Client   *http.Client
	Opt      Opt
}

var reProtocol = regexp.MustCompile("^https?://")

type Opt struct {
	Domain   string
	SkipPing bool
	Timeout  time.Duration
	NonSSL   bool
	Insecure bool
	Debug    bool
	Headers  map[string]string
}

func New(ctx context.Context, auth AuthConfig, opt Opt) (*registryClient, error) {
	transport := http.DefaultTransport
	return newFromTransport(ctx, auth, transport, opt)
}

func newFromTransport(ctx context.Context, auth AuthConfig, transport http.RoundTripper, opt Opt) (*registryClient, error) {
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
	registry := &registryClient{
		URL:    url,
		Domain: reProtocol.ReplaceAllString(url, ""),
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				//log.Info().Msg("There was a redirect!!!")
				//				return http.ErrUseLastResponse
				return nil
			},
			Timeout:   opt.Timeout,
			Transport: customTransport,
		},
		Username: auth.Username,
		Password: auth.Password,
	}
	// if registry.Pingable() && !opt.SkipPing {
	// 	if err := registry.Ping(ctx); err != nil {
	// 		return nil, err
	// 	}
	// }
	return registry, nil
}

func (r *registryClient) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *registryClient) getJSON(ctx context.Context, url string, response interface{}) (http.Header, error) {
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
	////log.Info().Msgf("registry.registry resp.Status=%s", resp.Status)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}
	return resp.Header, nil
}
