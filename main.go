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
	log.Info().Str("foo", "bar").Msg("Hello world")

}
func main() {
	fs := filestore.NewFileStore("/tmp/filestore")
	repoIndex := repo.NewRepoIndex()
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	//	repoRouter := r.PathPrefix("/repo").Subrouter()
	// Docker proxy
	dockerproxy.RegisterHandlers(r, fs, repoIndex)
	// Docker registry
	// Helm proxy
	// Go proxy
	// Local Docker registry
	// Local Go registry
	// Local Helm registry
	log.Print(apiRouter)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
	})
	log.Fatal().Err(http.ListenAndServe(":8081", r))
}
