package mocks

import (
	"context"

	"github.com/martencassel/binaryrepo"
	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type UploaderMock struct {
	mock.Mock
}

func NewUploaderMock() *UploaderMock {
	return &UploaderMock{}
}

func (m *UploaderMock) Get(ctx context.Context, name string) (*binaryrepo.User, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*binaryrepo.User), args.Error(1)
}

func (m *UploaderMock) AppendFile(uuid string, bytes []byte) error {
	args := m.Called(uuid, bytes)
	return args.Error(0)
}

func (m *UploaderMock) CreateUpload() (uuid.UUID, error) {
	args := m.Called()
	return args.Get(0).(uuid.UUID), args.Error(1)
}


func (m *UploaderMock) Exists(uuid string) bool {
	args := m.Called(uuid)
	return args.Bool(0)
}

func (m *UploaderMock) ReadUpload(uuid string) ([]byte, error) {
	args := m.Called(uuid)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *UploaderMock) Remove(uuid string) error {
	args := m.Called(uuid)
	return args.Error(0)
}

func (m *UploaderMock) WriteFile(uuid string, bytes []byte) error {
	args := m.Called(uuid, bytes)
	return args.Error(0)
}