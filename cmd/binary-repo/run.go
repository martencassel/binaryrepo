package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/martencassel/binaryrepo/pkg/api"
	"github.com/spf13/cobra"

	"github.com/rs/zerolog/log"
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

		// API router
		apiRouter := r.PathPrefix("/api").Subrouter()
		apiRouter.Use(loggingMiddleware)
		api.RegisterHandlers(apiRouter)

		srv := &http.Server{
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
			Handler:           r,
			Addr:              ":8081",
		}
		log.Info().Msgf("Listening on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	},
}
