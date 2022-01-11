package main

import "github.com/opencontainers/go-digest"

type FileStore interface {
	NewFileStore(basePath string) *FileStore
	Remove(digest digest.Digest) error
	Exists(digest digest.Digest) bool
	WriteFile(b []byte) (digest.Digest, error)
	ReadFile(digest digest.Digest) ([]byte, error)
	Size(digest digest.Digest) (int64, error)
}
