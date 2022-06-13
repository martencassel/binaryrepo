package dockerregistry

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/containers/image/manifest"
	"github.com/gorilla/mux"
	"github.com/opencontainers/go-digest"
	"github.com/rs/zerolog/log"
)

/*
	HEAD /v2/<name>/manifests/<reference>		reference = <digest> | <tag>
*/
func (registry *DockerRegistryHandler) ExistingManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.ExistingManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["namespace"]
	reference := vars["reference"]
	dgst, err := digest.Parse(reference)
	if err != nil {
		if !registry.ts.Exists(name, dgst.String()) {
			log.Printf("Could not find tag %s:%s", name, reference)
			msg := fmt.Sprintf("manifest for %s/%s not found", name, reference)
			writeErrorResponse(rw, http.StatusNotFound, ErrorDetails{
				Code:   fmt.Sprintf("%d", http.StatusNotFound),
				Message: msg,
				Details: msg,
			})
			return
		}
		digestStr, err := registry.ts.GetTag(name, reference)
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Found tag %s:%s with digest %s", name, reference, digestStr)
		dgst, err = digest.Parse(digestStr.String())
		if err != nil {
			log.Print(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if !registry.fs.Exists(dgst) {
		msg := fmt.Sprintf("manifest for %s/%s not found", name, reference)
		writeErrorResponse(rw, http.StatusNotFound, ErrorDetails{
			Code:   fmt.Sprintf("%d", http.StatusNotFound),
			Message: msg,
			Details: msg,
		})
		return
	}
	b , err := registry.fs.ReadFile(dgst)
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	contentLength := len(b)

	rw.Header().Set("Docker-Content-Digest", dgst.String())
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
}

/*
	GET /v2/<name>/manifests/<reference>		reference = <digest> | <tag>
*/
func (registry *DockerRegistryHandler) GetManifest(rw http.ResponseWriter, req *http.Request) {
	log.Info().Msgf("registry.GetManifest %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	name := vars["namespace"]
	reference := vars["reference"]
	log.Info().Msgf("Namespace: %s, Reference: %s", name, reference)
	var b []byte
	// Check if reference is a tag or digest
	dgst, err := digest.Parse(reference)
	if err == nil {
		// Reference is a digest
		if !registry.fs.Exists(dgst) {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		b, err = registry.fs.ReadFile(dgst)
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
		rw.Write([]byte(b))
		return
	} else {
		tag := reference
		exists := registry.ts.Exists(name, tag)
		if !exists {
			msg := fmt.Sprintf("manifest for %s/%s not found", name, tag)
			writeErrorResponse(rw, http.StatusNotFound, ErrorDetails{
				Code:   fmt.Sprintf("%d", http.StatusNotFound),
				Message: msg,
				Details: msg,
			})
			return
		}
		if exists {
			// Reference is a tag
			tagInfo, err := registry.ts.GetTag(name, reference)
			if err != nil {
				log.Print(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Info().Msgf("Found tag %s:%s with digest %s", name, reference, tagInfo.String())
			if !registry.fs.Exists(dgst) {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
			b, err = registry.fs.ReadFile(dgst)
			if err != nil {
				log.Print(err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			contentType := manifest.GuessMIMEType(b)
			rw.Header().Set("Content-Type", contentType)
			rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
			rw.Header().Set("Docker-Content-Digest", dgst.String())
			rw.WriteHeader(http.StatusOK)
			return
		}
	}
}


/*
	PUT /v2/<name>/manifests/<reference>		reference = <digest> | <tag>
*/
func (registry *DockerRegistryHandler) UploadManifest(rw http.ResponseWriter, req *http.Request) {
	var isTag bool
	log.Info().Msgf("registry.hasManifestHandler %s %s", req.Method, req.URL.Path)
	vars := mux.Vars(req)
	namespace := vars["namespace"]
	reference := vars["reference"]
	digest, err:= digest.Parse(reference)
	if err != nil {
		isTag = true
	} else {
		isTag = false
	}
	log.Info().Msgf("Namespace: %s, Digest: %s", namespace, digest)


	if req.Body == http.NoBody || req.ContentLength == 0 {
		log.Info().Msg("no body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_digest, err := registry.fs.WriteFile(body);
	if err != nil {
		log.Print(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isTag {
		registry.ts.WriteTag(namespace, reference, _digest)
	}

	location := fmt.Sprintf("/v2/%s/manifests/%s", namespace, _digest.String())
	rw.Header().Set("Docker-Content-Digest", _digest.String())
	rw.Header().Set("docker-distribution-api-version", "registry/2.0")
	rw.Header().Set("Content-Length", "0")
	rw.Header().Set("Location", location)
	rw.WriteHeader(http.StatusCreated)
}

/*
	DELETE /v2/<name>/manifests/<reference>		reference = <digest> | <tag>
*/
func (registry *DockerRegistryHandler) DeleteManifest(rw http.ResponseWriter, req *http.Request) {
}