package cmd

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	dockerproxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
)

func Execute() {
	fs := filestore.NewFileStore("/tmp/filestore")
	repoIndex := repo.NewRepoIndex()
	repoIndex.AddRepo(repo.Repo{
		ID:      1,
		Name:    "docker-remote",
		Type:    repo.Remote,
		PkgType: repo.Docker,
		URL:     "https://registry-1.docker.io",
	})
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	dockerproxy.RegisterHandlers(r, fs, repoIndex)
	log.Print(apiRouter)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
	})
	srv := &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
		Addr:              ":8081",
	}
	log.Logger.Fatal().Err(srv.ListenAndServe())
}
