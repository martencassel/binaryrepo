package binaryrepo

import (
	"context"
	"time"

	"github.com/opencontainers/go-digest"
)

type MetadataInfo struct {
	Digest digest.Digest
	Repository string
	Path string
	Etag string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MetadataIndexer interface {
	GetByDigest(ctx context.Context, digest digest.Digest) (*MetadataInfo, error)
	GetByRepoPath(ctx context.Context, repo string, path string) (*MetadataInfo, error)
	Put(ctx context.Context, digest digest.Digest, path string, repo string)
}
