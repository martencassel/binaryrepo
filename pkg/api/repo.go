package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/martencassel/binaryrepo"
	"github.com/rs/zerolog/log"
)

func (api *ApiHandler) repoCreate(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("repoCreate: %s %s", req.Method, req.URL.Path)
	var repo binaryrepo.Repo
    err := json.NewDecoder(req.Body).Decode(&repo)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }
	log.Info().Msgf("%v", repo)
	ctx := context.Background()
	err = api.rs.Create(ctx, &repo)
	log.Print(err)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(`{"message": "repo created"}`))
 }

func (api *ApiHandler) repoList(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("repoList: %s %s", req.Method, req.URL.Path)
	list, err :=  api.rs.List(context.Background())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(list)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
	rw.Write(data)
}