package main

import (
	"net/http"
	"os"
	"time"

	dockerproxy "github.com/martencassel/binaryrepo/pkg/docker/proxy"
	dockerregistry "github.com/martencassel/binaryrepo/pkg/docker/registry"
	dockerrouter "github.com/martencassel/binaryrepo/pkg/docker/router"
	"github.com/martencassel/binaryrepo/pkg/docker/uploader"

	"github.com/gorilla/mux"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	"github.com/martencassel/binaryrepo/pkg/repo"
	version "github.com/martencassel/binaryrepo/pkg/util/version"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		log.Info().Msgf("loggingMiddleware: %s", req.RequestURI)
		next.ServeHTTP(rw, req)
	})
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a binaryrepo server",
	Long:  `Start a binaryrepo server`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Str("version", version.AppVersion()).Msg("running server")
		fs := filestore.NewFileStore("/tmp/filestore")
		uploader := uploader.NewUploadManager("/tmp/uploads")
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
		repoIndex.AddRepo(repo.Repo{
			ID:      2,
			Name:    "docker-local",
			Type:    repo.Local,
			PkgType: repo.Docker,
		})
		repoIndex.AddRepo(repo.Repo{
			ID: 3,
			Name: "docker-group",
			Group: []string{
				"docker-local", "docker-remote",
			},
		})
		r := mux.NewRouter()
		dockerProxy := dockerproxy.NewProxyAppWithOptions(fs, repoIndex)
		dockerRegistry := dockerregistry.NewDockerRegistry(fs, repoIndex, uploader)

		dockerRouter := dockerrouter.NewDockerRouter(dockerProxy, dockerRegistry, repoIndex)
		dockerRouter.RegisterHandlers(r)
		r.Use(loggingMiddleware)
		r.PathPrefix("/").HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			log.Info().Msgf("not-implemented %s %s", req.Method, req.URL)
			//			vars := mux.Vars(r)
			//			repoName := vars["repo-name"]
		})
		/*
			If the timeout are to low, DownloadLayer handler will failed with an error such as

				Error copying response: readfrom tcp [::1]:8081->[::1]:54774: write tcp [::1]:8081->[::1]:54774: i/o timeout

			For big layer blobs as in the postgres image

				docker-remote.example.com /v2/postgres/blobs/sha256:794976979956b97dc86e3b99fc0cdcd6385113969574152ba4a6218431f542e9

			This may happen
		*/
		srv := &http.Server{
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			Handler:           r,
			Addr:              ":8081",
		}
		////log.Info().Msgf("Listening on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	},
}
