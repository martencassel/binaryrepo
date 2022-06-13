package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/api"
	"github.com/martencassel/binaryrepo/pkg/fakes"
	filestore "github.com/martencassel/binaryrepo/pkg/filestore/fs"
	postgres "github.com/martencassel/binaryrepo/pkg/postgres"
	"github.com/martencassel/binaryrepo/pkg/repostore"
	userstore "github.com/martencassel/binaryrepo/pkg/userstore"
	"github.com/spf13/cobra"

	"github.com/rs/zerolog/log"

	dockerregistry "github.com/martencassel/binaryrepo/pkg/docker/registry"
	registry "github.com/martencassel/binaryrepo/pkg/docker/registry"
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
		// Base router
		r := mux.NewRouter()

		// Database
		ctx := context.Background()
		db, err := postgres.Open(ctx, binaryrepo.Config{
			Host: "binaryrepo-postgres",
			Port: "5432",
			User: "postgres",
			Password: "postgres",
			DBName: "binaryrepo",
			MaxOpenConnections: 10,
			MaxIdleConnections: 5,
			CreateDB: true,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open database")
		}
		log.Info().Msg("Database connection established")
		defer db.Close()

		// Dependencies
		userStore :=  userstore.NewUserStore(db)
		repoStore := repostore.NewRepoStore(db)
		tagStore := fakes.NewTagStore()
		//nodeStore := nodestore.NewNodeStore(db)
		fs := filestore.NewFileStore("/tmp/filestore")

		index := fakes.NewIndexer()
//		uploader := uploader.NewUploader(index)
		fakeUploader := fakes.NewUploader(index)

		// API handler
		apiRouter := r.PathPrefix("/api").Subrouter()
		apiRouter.Use(loggingMiddleware)
		apiAuth := api.NewBasicAuthMiddleware(userStore)
		apiRouter.Use(apiAuth.Middleware)
		api := api.NewApiHandler(repoStore, userStore, fs)
		api.RegisterHandlers(apiRouter)

		// Docker Registry
		dockerRouter := r.PathPrefix("/").Subrouter()
		dockerRouter.Use(loggingMiddleware)
		registryConfig := &registry.RegistryConfig{
			RepoName: "docker-local",
			Port: "443",
			DockerDomain: "binaryrepo.local",
			DockerPort: "443",
		}
		mw := dockerregistry.NewRegistryV2AuthMiddleware( registryConfig, & userStore)
		dockerRouter.Use(mw.Middleware)
		registry := registry.NewDockerRegistryHandler(repoStore, fs, tagStore, fakeUploader)
		registry.RegisterHandlers(dockerRouter)

		// Server
		srv := &http.Server{
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			Handler:           r,
			Addr:              "0.0.0.0:8081",
		}
		log.Info().Msgf("Listening on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	},
}
