package mw

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info().Msg(r.RequestURI)
			val, ok := r.Header[http.CanonicalHeaderKey("X-API-KEY")]
			if ok {
					log.Info().Msgf("X-API-KEY key header is present %s\n", val)
			} else {
					log.Info().Msg("X-API-KEY key header is not present")
					next.ServeHTTP(w, r)
					return
			}
			// Check against all API key's
			apiKey := r.Header.Get("X-API-KEY")
			if apiKey == "ABcdEF" {
					next.ServeHTTP(w, r)
					return
			}
			log.Error().Msg("Unauthorized API key")
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized X-API-KEY for user", http.StatusUnauthorized)
	})
}
