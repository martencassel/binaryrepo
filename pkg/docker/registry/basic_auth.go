package dockerregistry

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/martencassel/binaryrepo"
	"github.com/rs/zerolog/log"
)

type registryV2AuthMiddleware struct {
	reponame string
	port string
	docker_domain string
	docker_port string
	userStore binaryrepo.UserStore
}

type RegistryConfig struct {
	RepoName string
	Port string
	DockerDomain string
	DockerPort string
}

func NewRegistryV2AuthMiddleware(config *RegistryConfig, userStore *binaryrepo.UserStore) *registryV2AuthMiddleware {
	return &registryV2AuthMiddleware{
		reponame: config.RepoName,
		port: config.Port,
		docker_domain: config.DockerDomain,
		docker_port: config.DockerPort,
		userStore: *userStore,
	}
}

func (mw *registryV2AuthMiddleware) getAuthInfo(docker_domain, docker_port, reponame, path string) string {
	realm := fmt.Sprintf("%s.%s/token", reponame, docker_domain)
	serviceString := docker_domain
	scopeString := fmt.Sprintf("repository:%s:pull", reponame)
	authValue := fmt.Sprintf("Basic realm=\"https://%s\",service=\"%s\",scope=\"%s\"", realm, serviceString, scopeString)
	log.Info().Msgf("auth value: %s\n", authValue)
	return authValue

}

func (mw *registryV2AuthMiddleware) Middleware(next http.Handler)  http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msgf("registryV2AuthMiddleware: %s %s %s", r.Method, r.URL.Path, path.Base(r.URL.Path))
		username, password, ok := r.BasicAuth()
		if ok {
			log.Info().Msgf("registryV2AuthMiddleware: BasicAuth, username: %s, password: %s", username, password)
			ctx := context.Background()
			found, err := mw.userStore.LookupUser(ctx, username, password);
			if err != nil {
				log.Error().Err(err).Msg("registryV2AuthMiddleware: LookupUser failed")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if found == true {
				log.Info().Msgf("registryV2AuthMiddleware: LookupUser success")
				next.ServeHTTP(w, r)
				return
			}
		}
		log.Info().Msgf("registryV2AuthMiddleware: BasicAuth failed")
		authInfo := mw.getAuthInfo(mw.docker_domain, mw.docker_port, mw.reponame, r.URL.Path)
		log.Info().Msgf("Auth Info: %s", authInfo)
		data := (`{"errors":[{"code": "UNAUTHORIZED","message": "authentication required","detail":null}]}`)
		w.Header().Set("WWW-Authenticate", authInfo)
		w.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(data))
	})
}

func Middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // do stuff
        h.ServeHTTP(w, r)
    })
}