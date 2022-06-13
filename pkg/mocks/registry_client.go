package mocks

import (
	"context"
	"io"
	"net/http"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/martencassel/binaryrepo/pkg/docker/client"
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/mock"
)

type RegistryClientMock struct {
	mock.Mock
}

func NewRegistryClientMock() *RegistryClientMock {
	return &RegistryClientMock{}
}

func (m *RegistryClientMock) SetConfig(url string, domain string, config *client.AuthConfig) {
	m.Called(url, domain, config)
}

func (m *RegistryClientMock) DownloadLayer(ctx context.Context, repository string, digest digest.Digest) (io.ReadCloser, *http.Response, error) {
	args := m.Called(ctx, repository, digest)
	var r io.ReadCloser
	var rsp *http.Response
	if args.Get(0) != nil {
		r = args.Get(0).(io.ReadCloser)
		rsp = args.Get(1).(*http.Response)
	}
	return r, rsp, args.Error(2)
}

func (m *RegistryClientMock) HasLayer(ctx context.Context, repository string, digest digest.Digest) (bool, *http.Response, error) {
	args := m.Called(ctx, repository, digest)
	return args.Get(0).(bool), args.Get(1).(*http.Response), args.Error(2)
}

func (m *RegistryClientMock) Pingable() bool {
	args := m.Called()
	return args.Get(0).(bool)
}

func (m *RegistryClientMock) Ping(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Get(0).(bool)
}

func (m *RegistryClientMock) Manifest(ctx context.Context, repository, ref string) (distribution.Manifest, error) {
	args := m.Called(ctx, repository, ref)
	return args.Get(0).(distribution.Manifest), args.Error(1)
}

func (m *RegistryClientMock) ManifestList(ctx context.Context, repository, ref string) (manifestlist.ManifestList, error) {
	args := m.Called(ctx, repository, ref)
	return args.Get(0).(manifestlist.ManifestList), args.Error(1)
}

func (m *RegistryClientMock) ManifestV2(ctx context.Context, repository, ref string) (schema2.Manifest, error) {
	args := m.Called(ctx, repository, ref)
	return args.Get(0).(schema2.Manifest), args.Error(1)
}

func (m *RegistryClientMock) ManifestV1(ctx context.Context, repository, ref string) (schema1.SignedManifest, error) {
	args := m.Called(ctx, repository, ref)
	return args.Get(0).(schema1.SignedManifest), args.Error(1)
}

func (m *RegistryClientMock) Token(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.Get(0).(string), args.Error(1)
}
