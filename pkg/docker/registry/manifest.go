package registry

import (
	"fmt"
	"io"
	"net/http"

	"github.com/containers/image/manifest"
	"github.com/gorilla/mux"
	digest "github.com/opencontainers/go-digest"
	log "github.com/rs/zerolog/log"
)

func (registry *DockerRegistry) HasManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.hasManifestHandler %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]
	reference := vars["reference"]
	dgst, err := digest.Parse(reference)
	if err != nil {
		if !registry.tagstore.TagExists(name, reference) {
			log.Printf("Could not find tag %s:%s", name, reference)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		digestStr, err := registry.tagstore.ReadTag(name, reference)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Found tag %s:%s with digest %s", name, reference, digestStr)
		dgst, err = digest.Parse(digestStr)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if !registry.fs.Exists(dgst) {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := registry.fs.ReadFile(dgst)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Docker-Content-Digest", dgst.String())
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	rw.WriteHeader(http.StatusOK)

}

func (registry *DockerRegistry) GetManifestHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.getManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]
	reference := vars["reference"]
	dgst, err := digest.Parse(reference)
	if err != nil {
		if !registry.tagstore.TagExists(name, reference) {
			log.Printf("Could not find tag %s:%s", name, reference)
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		digestStr, err := registry.tagstore.ReadTag(name, reference)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Found tag %s:%s with digest %s", name, reference, digestStr)
		dgst, err = digest.Parse(digestStr)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	b, err := registry.fs.ReadFile(dgst)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	contentType := manifest.GuessMIMEType(b)
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Docker-Content-Digest", dgst.String())
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(b)
	if err != nil {
		log.Print(err)
		return
	}
}

// PathPutManifest URL.
const PathPutManifest = "/repo/{repo-name}/v2/{name}/manifests/{reference}"

// Put the manifest identified by name and reference where reference can be a tag or digest.
func (registry *DockerRegistry) PutManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.putManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["name"]
	reference := vars["reference"]
	b, err := io.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	digest, err := registry.fs.WriteFile(b)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	exists := registry.fs.Exists(digest)
	if !exists {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	// Write manifest file to storage
	digest, err = registry.fs.WriteFile(b)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Docker-Content-Digest", digest.String())
	rw.Header().Set("Location", "/v2/"+name+"/manifests/"+reference)
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
	rw.WriteHeader(http.StatusCreated)
}
