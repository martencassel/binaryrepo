package indexer

import (
	"context"

	"github.com/martencassel/binaryrepo"
	"github.com/opencontainers/go-digest"
)

type indexMetadata struct {
	ns *binaryrepo.NodeStore
}

func NewIndexer(ns *binaryrepo.NodeStore) binaryrepo.MetadataIndexer {
	return &indexMetadata{
		ns: ns,
	}
}


func (indexer *indexMetadata) Put(ctx context.Context, digest digest.Digest, path string, repo string) {
	nodeStore := *indexer.ns
	node := &binaryrepo.Node{
		RepoName: repo,
		Path: path,
		Checksum: digest.String(),
	}
	nodeStore.UpdateNode(ctx, node)
}

func (indexer *indexMetadata) GetByDigest(ctx context.Context, digest digest.Digest) (*binaryrepo.MetadataInfo, error) {
	nodeStore := *indexer.ns
	node, err := nodeStore.GetNodeByDigest(ctx, digest)
	if err != nil {
		return nil, err
	}
	info := &binaryrepo.MetadataInfo{}
	info.Digest = digest.Algorithm().FromString(node.Checksum)
	info.Path = node.Path
	info.Repository = node.RepoName
	return info, nil
}

func (indexer *indexMetadata) GetByRepoPath(ctx context.Context, repo string, path string) (*binaryrepo.MetadataInfo, error) {
	return nil, nil
}

