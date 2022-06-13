package local

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func AlpineGetHandler() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Info().Msg(r.RequestURI)
		log.Info().Msg("Alpine local get handler")
		vars := mux.Vars(r)
		log.Info().Msg("Handle GET")
		log.Info().Msgf("%s %s %s %s", vars["branch"], vars["repo"], vars["arch"], vars["filename"])
	})
}