package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/martencassel/binaryrepo"
	"github.com/rs/zerolog/log"
)

type basicAuthMiddleware struct {
	userStore binaryrepo.UserStore
}

func NewBasicAuthMiddleware(userStore binaryrepo.UserStore) * basicAuthMiddleware {
	return &basicAuthMiddleware{
		userStore: userStore,
	}
}

func (mw *basicAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			log.Info().Msgf("API basicAuthMiddleware: BasicAuth, username: %s, password: %s", username, password)
			ctx := context.Background()
			found, err := mw.userStore.LookupUser(ctx, username, password);
			if err != nil {
				log.Error().Err(err).Msg("API basicAuthMiddleware: LookupUser failed")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if found == true {
				log.Info().Msgf("API basicAuthMiddleware: LookupUser success")
				next.ServeHTTP(w, r)
				return
			}
		}
		log.Info().Msgf("basicAuthMiddleware: BasicAuth failed")
		data := (`{"errors":[{"code": "UNAUTHORIZED","message": "authentication required","detail":null}]}`)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(data))
	})
}