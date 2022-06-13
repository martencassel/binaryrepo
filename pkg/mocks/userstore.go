package mocks

import (
	"context"

	"github.com/martencassel/binaryrepo"
	"github.com/stretchr/testify/mock"
)

type UserStoreMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserStoreMock {
	return &UserStoreMock{}
}

func (m *UserStoreMock) GetUser(ctx context.Context, name string) (*binaryrepo.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*binaryrepo.User), args.Error(1)
}

func (m *UserStoreMock) LookupUser(ctx context.Context, username string, password string) (bool, error) {
	args := m.Called(ctx, username, password)
	return args.Bool(0), args.Error(1)
}

func (m *UserStoreMock) Create(ctx context.Context, repo *binaryrepo.User) error {
	args := m.Called(ctx, repo)
	return args.Get(0).(error)
}

func (m *UserStoreMock) List(ctx context.Context) ([]binaryrepo.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]binaryrepo.User), args.Error(1)
}
