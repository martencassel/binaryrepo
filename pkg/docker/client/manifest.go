package client

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/rs/zerolog/log"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
)

var (
	// ErrUnexpectedSchemaVersion a specific schema version was requested, but was not returned
	ErrUnexpectedSchemaVersion = errors.New("recieved a different schema version than expected")
)

// Manifest returns the manifest for a specific repository:tag.
func (r *Registry) Manifest(ctx context.Context, repository, ref string) (distribution.Manifest, error) {
	uri := r.url("/v2/%s/manifests/%s", repository, ref)
	log.Printf("registry.manifests uri=%s repository=%s ref=%s", uri, repository, ref)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", schema2.MediaTypeManifest)
	resp, err := r.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	//defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("registry.manifests resp.Status=%s, body=%s", resp.Status, body)
	m, _, err := distribution.UnmarshalManifest(resp.Header.Get("Content-Type"), body)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ManifestList gets the registry v2 manifest list.
func (r *Registry) ManifestList(ctx context.Context, repository, ref string) (manifestlist.ManifestList, error) {
	uri := r.url("/v2/%s/manifests/%s", repository, ref)
	log.Printf("registry.manifests uri=%s repository=%s ref=%s", uri, repository, ref)
	var m manifestlist.ManifestList
	if _, err := r.getJSON(ctx, uri, &m); err != nil {
		//log.Info().Msg(err.Error())
		////log.Info().Msgf("registry.manifests response=%v", m)
		return m, err
	}
	return m, nil
}

// ManifestV2 gets the registry v2 manifest.
func (r *Registry) ManifestV2(ctx context.Context, repository, ref string) (schema2.Manifest, error) {
	uri := r.url("/v2/%s/manifests/%s", repository, ref)
	////log.Info().Msgf("registry.manifests uri=%s repository=%s ref=%s", uri, repository, ref)
	var m schema2.Manifest
	if _, err := r.getJSON(ctx, uri, &m); err != nil {
		////log.Info().Msgf("registry.manifests response=%v", m)
		return m, err
	}
	if m.Versioned.SchemaVersion != 2 {
		return m, ErrUnexpectedSchemaVersion
	}
	return m, nil
}

// ManifestV1 gets the registry v1 manifest.
func (r *Registry) ManifestV1(ctx context.Context, repository, ref string) (schema1.SignedManifest, error) {
	uri := r.url("/v2/%s/manifests/%s", repository, ref)
	////log.Info().Msgf("registry.manifests uri=%s repository=%s ref=%s", uri, repository, ref)
	var m schema1.SignedManifest
	if _, err := r.getJSON(ctx, uri, &m); err != nil {
		////log.Info().Msgf("registry.manifests response=%v", m)
		return m, err
	}
	if m.Versioned.SchemaVersion != 1 {
		return m, ErrUnexpectedSchemaVersion
	}
	return m, nil
}
