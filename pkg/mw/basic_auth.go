package mw

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/rs/zerolog/log"
)

func basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Info().Msg(r.RequestURI)
		username, password, ok := r.BasicAuth()
		if ok {
			  // Check username and password
			  usernameHash := sha256.Sum256([]byte(username))
			  passwordHash := sha256.Sum256([]byte(password))
			  expectedUsernameHash := sha256.Sum256([]byte("admin"))
			  expectedPasswordHash := sha256.Sum256([]byte("P@ssw0rd"))
			  usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			  passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			  if usernameMatch && passwordMatch {
					  next.ServeHTTP(w, r)
					  return
			  }

			  // Check username and personal API key
			  expectedAPIKeyHash := sha256.Sum256([]byte("ABcdEF"))
			  passwordMatch = (subtle.ConstantTimeCompare(passwordHash[:], expectedAPIKeyHash[:]) == 1)
			  if usernameMatch && passwordMatch {
					  next.ServeHTTP(w, r)
					  return
			  }

		}
		log.Error().Msg("Unauthorized basic auth")
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}