package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	dockerproxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	version "github.com/martencassel/binaryrepo/pkg/util/version"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a binaryrepo server",
	Long:  `Start a binaryrepo server`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Str("version", version.AppVersion()).Msg("running server")
		fs := filestore.NewFileStore("/tmp/filestore")
		repoIndex := repo.NewRepoIndex()
		hubUser := os.Getenv("DOCKERHUB_USERNAME")
		hubPass := os.Getenv("DOCKERHUB_PASSWORD")
		repoIndex.AddRepo(repo.Repo{
			ID:       1,
			Name:     "docker-remote",
			Type:     repo.Remote,
			PkgType:  repo.Docker,
			URL:      "https://registry-1.docker.io",
			Username: hubUser,
			Password: hubPass,
		})
		r := mux.NewRouter()
		//r.PathPrefix("/api").Subrouter()
		dockerproxy.RegisterHandlers(r, fs, repoIndex)
		r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL)
		})
		srv := &http.Server{
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			Handler:           r,
			Addr:              ":8081",
		}
		log.Info().Msgf("Listening on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	},
}
