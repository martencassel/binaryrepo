package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func CreateMocks() (*mocks.FileStoreMock, *mocks.UserStoreMock,
	*mocks.RepoStoreMock) {
	fileStoreMock := mocks.NewFileStoreMock()
	userStoreMock := mocks.NewUserServiceMock()
	repoStoreMock := mocks.NewRepoStoreMock()
	return fileStoreMock, userStoreMock, repoStoreMock
}

func TestApi(t *testing.T) {
	t.Run("Can create a repo", func(t *testing.T) {
		// Arrange
		repoName := "docker-local"
		repo := binaryrepo.Repo {
			Name: repoName,
			Repotype: "local",
			Pkgtype: "docker",
		}
		fs, us, rs := CreateMocks()
		rs.On("Create", mock.Anything, &repo).Return(nil)
		apiHandler := NewApiHandler(rs, us, fs)
		res := httptest.NewRecorder()
		b, err := json.Marshal(repo)
		assert.NoError(t, err)
		body := bytes.NewBuffer(b)
		req, _ := http.NewRequest(http.MethodPost, "", body)
		// Act
		apiHandler.repoCreate(res, req)
		// Assert
		assert.NoError(t, err)
		assert.True(t, res.Code == http.StatusCreated)
	})
}