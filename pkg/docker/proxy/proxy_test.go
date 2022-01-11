package dockerproxy

import (
	"net/http"
	"net/http/httptest"
)

func DockerHubV2ResponseStub() *httptest.Server {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("docker-distribution-api-version", "registry/2.0")
		w.Header().Add("www-authenticate", "Bearer realm=\"https://auth.docker.io/token\",service=\"registry.docker.io\"")
	}))
	return svr
}
