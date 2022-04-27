package api

import (
	"encoding/json"
	"net/http"

	binaryrepo "github.com/martencassel/binaryrepo"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/repo", createRepo).Methods(http.MethodPost)
	r.HandleFunc("/repo", listRepos).Methods(http.MethodGet)
}

func createRepo(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("createRepo: %s %s", req.Method, req.URL.Path)

	var repo binaryrepo.Repo

    err := json.NewDecoder(req.Body).Decode(&repo)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusBadRequest)
        return
    }
	log.Info().Msgf("%v", repo)
 }

func listRepos(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("listRepos: %s %s", req.Method, req.URL.Path)
}