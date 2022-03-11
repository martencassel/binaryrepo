package binaryrepo

import "github.com/opencontainers/go-digest"

// Filestore is an interface for a service that stores binary blobs.
type Filestore interface {
	// Exists checks if a binary blob exists in the filestore.
	Exists(digest digest.Digest) bool

	// Write a file to the filestore.
	WriteFile(b []byte) (digest.Digest, error)

	// Read a file from the filestore.
	ReadFile(digest digest.Digest) ([]byte, error)

	// Remove removes a binary blob from the filestore.
	Remove(digest digest.Digest) error
}