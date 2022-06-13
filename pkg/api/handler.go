package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (api *ApiHandler) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/repo", api.repoCreate).Methods(http.MethodPost)
	r.HandleFunc("/repo", api.repoList).Methods(http.MethodGet)
}
