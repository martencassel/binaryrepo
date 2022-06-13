package dockerregistry

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)


type TagList struct {
	Name string
	Tags []string
}

/*
	GET /v2/<name>/tags/list
*/
func (registry *DockerRegistryHandler) ListImageTags(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.ListImageTags %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["namespace"]
	tags, err := registry.ts.GetTags(name)
	tagList := TagList{
		Name: name,
		Tags: tags,
	}
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	b, err := json.Marshal(tagList)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(b))
}
