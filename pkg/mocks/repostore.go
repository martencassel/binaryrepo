package mocks

import (
	"context"

	"github.com/martencassel/binaryrepo"
	"github.com/stretchr/testify/mock"
	//	postgres "github.com/martencassel/binaryrepo/pkg/postgres"
)

type RepoStoreMock struct {
	mock.Mock
}

func NewRepoStoreMock() *RepoStoreMock {
	return &RepoStoreMock{}
}

func (m *RepoStoreMock) Get(ctx context.Context, name string) (*binaryrepo.Repo, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*binaryrepo.Repo), args.Error(1)
}

func (m *RepoStoreMock) Create(ctx context.Context, repo *binaryrepo.Repo) error {
	args := m.Called(ctx, repo)
	return args.Error(0)
}

func (m *RepoStoreMock) List(ctx context.Context) ([]binaryrepo.Repo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]binaryrepo.Repo), args.Error(1)
}

func (m *RepoStoreMock) Exists(ctx context.Context, repoName string) (bool, error) {
	args := m.Called(ctx)
	return args.Get(0).(bool), args.Error(1)
}
