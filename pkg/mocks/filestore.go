package mocks

import (
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/mock"
)

type FileStoreMock struct {
	mock.Mock
}

func NewFileStoreMock() *FileStoreMock {
	return &FileStoreMock{}
}


func (m *FileStoreMock) Remove(digest digest.Digest) error {
	args := m.Called(digest)
	return args.Get(0).(error)
}

func (m *FileStoreMock) Exists(digest digest.Digest) bool {
	args := m.Called(digest)
	return args.Get(0).(bool)
}


func (m *FileStoreMock) GetBasePath() string {
	args := m.Called()
	return args.String(0)
}

func (m *FileStoreMock) WriteFile(b []byte) (digest.Digest, error) {
	args := m.Called(b)
	return args.Get(0).(digest.Digest), args.Error(1)
}

func (m *FileStoreMock) ReadFile(digest digest.Digest) ([]byte, error) {
	args := m.Called(digest)
	return args.Get(0).([]byte), args.Error(1)
}

