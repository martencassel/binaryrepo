package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	dockerproxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting.....")

}
func main() {
	fs := filestore.NewFileStore("/tmp/filestore")
	repoIndex := repo.NewRepoIndex()
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	dockerproxy.RegisterHandlers(r, fs, repoIndex)
	log.Print(apiRouter)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
	})
	log.Fatal().Err(http.ListenAndServe(":8081", r))
}
