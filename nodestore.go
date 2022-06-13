package binaryrepo

import (
	"context"

	"github.com/opencontainers/go-digest"
)

type Node struct {
	ID 	  int64
	RepoName string
	RepoID  int64
	NodeName string
	Path string
	UpstreamUrl string
	ETag string
	Checksum string
	IsFolder bool
	ParentID int64
}

type NodeStore interface {
	CreateNode(ctx context.Context, node Node) error
	GetNodeByDigest(ctx context.Context, digest digest.Digest) (*Node, error)
	GetNode(ctx context.Context, repo string, path string) (*Node, error)
	UpdateNode(ctx context.Context, node *Node) error
}