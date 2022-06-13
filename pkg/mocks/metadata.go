package mocks

import (
	"context"

	metadata "github.com/martencassel/binaryrepo"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/mock"
)

type MetadataStoreMock struct {
	mock.Mock
}

func NewMetadataStoreMock() *MetadataStoreMock {
	return &MetadataStoreMock{}
}

func (m *MetadataStoreMock) GetByDigest(ctx context.Context, digest digest.Digest) (*metadata.MetadataInfo, error) {
	args := m.Called(ctx, digest)
	return args.Get(0).(*metadata.MetadataInfo), args.Error(1)
}

func (m *MetadataStoreMock) GetByRepoPath(ctx context.Context, repo string, path string) (*metadata.MetadataInfo, error) {
	args := m.Called(ctx, repo, path)
	return args.Get(0).(*metadata.MetadataInfo), args.Error(1)
}

func (m *MetadataStoreMock) Put(ctx context.Context, digest digest.Digest, path string, repo string) {
	m.Called(ctx, digest, path, repo)
}
