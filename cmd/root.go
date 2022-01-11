package cmd

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	dockerproxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

func Execute() {
	fs := filestore.NewFileStore("/tmp/filestore")
	repoIndex := repo.NewRepoIndex()
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	dockerproxy.RegisterHandlers(r, fs, repoIndex)
	log.Print(apiRouter)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
	})
	log.Logger.Fatal().Err(http.ListenAndServe(":8081", r))
}
