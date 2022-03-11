package binaryrepo

import "github.com/opencontainers/go-digest"

/*
	The tagstore support managing tags for the following registry operations:

	GET /<v2/<name>/manifests/<reference>			reference := <tag> | <digest>
	HEAD /v2/<name>/manifests/<reference>			reference := <tag> | <digest>
	PUT /v2/<name>/manifests/<reference>			reference := <tag> | <digest>
	DELETE /v2/<name>/manifests/<reference>			reference := <tag> | <digest>
*/

// Tag store is an interface for a service that manages tags for docker images
type TagStore interface {

	// Exists checks if a tag exists in the tag store
	Exists(repo string, digest digest.Digest) bool

	// Returns the list of tags for a given image repository
	GetTags(repo string) ([]string, error)

	// Writes a tag to the tag store.
	WriteTag(repo, tag, digest string) error

	// Get a tag from the tag store, a tag is a reference to a digest.
	GetTag(repo, tag string) (digest.Digest, error)
}

