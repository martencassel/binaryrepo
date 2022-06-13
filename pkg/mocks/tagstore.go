package mocks

import (
	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/mock"
	//	postgres "github.com/martencassel/binaryrepo/pkg/postgres"
)

type TagStoreMock struct {
	mock.Mock
}

func NewTagStoreMock() *TagStoreMock {
	return &TagStoreMock{}
}

func (t *TagStoreMock) Exists(repo string, tag string) bool {
	args := t.Called(repo, tag)
	return args.Bool(0)
}

func (t *TagStoreMock) GetTags(repo string) ([]string, error) {
	args := t.Called(repo)
	return args.Get(0).([]string), args.Error(1)
}

func (t *TagStoreMock) 	WriteTag(repo, tag string, digest digest.Digest) error {
	args := t.Called(repo, tag, digest)
	return args.Error(0)
}

func (t *TagStoreMock) GetTag(repo, tag string) (digest.Digest, error) {
	args := t.Called(repo, tag)
	return args.Get(0).(digest.Digest), args.Error(1)
}