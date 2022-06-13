package binaryrepo

import "github.com/opencontainers/go-digest"

// Tag store is an interface for a service that manages tags for docker images
type TagStore interface {

	// Exists checks if a tag exists in the tag store
	Exists(repo string, tag string) bool

	// Returns the list of tags for a given image repository
	GetTags(repo string) ([]string, error)

	// Writes a tag to the tag store.

	WriteTag(repo string, tag string, digest digest.Digest) error

	// Get a tag from the tag store, a tag is a reference to a digest.
	GetTag(repo, tag string) (digest.Digest, error)
}

