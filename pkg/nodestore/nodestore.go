package nodestore

import (
	"context"

	binaryrepo "github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/postgres"
	"github.com/opencontainers/go-digest"
)

type nodeStore struct {
	db *postgres.Database
}

func NewNodeStore(db *postgres.Database) binaryrepo.NodeStore {
	return &nodeStore{
		db: db,
	}
}

func (n *nodeStore) CreateNode(ctx context.Context, node binaryrepo.Node) error {
	txn, err := n.db.Begin(ctx)
	if err != nil {
		return err
	}
	err = txn.CreateNode(node)
	if err != nil {
		return err
	}
	return txn.Commit(ctx)
}

func (n *nodeStore) GetNode(ctx context.Context, repo string, path string) (*binaryrepo.Node, error) {
	txn, err := n.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	node, err := txn.GetNode(repo, path)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (n *nodeStore) GetNodeByDigest(ctx context.Context, digest digest.Digest) (*binaryrepo.Node, error) {
	txn, err := n.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	_, err = txn.GetNodeByDigest(digest)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (n *nodeStore) UpdateNode(ctx context.Context, node *binaryrepo.Node) (error) {
	txn, err := n.db.Begin(ctx)
	if err != nil {
		return err
	}
	err = txn.UpdateNode(*node)
	if err != nil {
		return err
	}
	return txn.Commit(ctx)
}