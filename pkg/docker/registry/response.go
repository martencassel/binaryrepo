package dockerregistry

import (
	"encoding/json"
	"net/http"
)

type ErrorDetails struct {
	Code string
	Message string
	Details string
}

func writeErrorResponse(rw http.ResponseWriter, status int, manifestError ErrorDetails) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(manifestError)
}
