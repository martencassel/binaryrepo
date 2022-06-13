package tagstore

import (
	binaryrepo "github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/postgres"
	"github.com/opencontainers/go-digest"
)

type tagStore struct {
	db *postgres.Database
}

func NewTagStore(db *postgres.Database) binaryrepo.TagStore {
	return &tagStore{
		db: db,
	}
}

func (ts *tagStore) Exists(repo string, tag string) bool {
	return false
}

func (ts *tagStore) GetTags(repo string) ([]string, error) {
	return nil, nil
}

func (t *tagStore) WriteTag(repo string, tag string, digest digest.Digest) error {
	return nil
}

func (ts *tagStore) GetTag(repo, tag string) (digest.Digest, error) {
	return digest.Digest(""), nil
}
