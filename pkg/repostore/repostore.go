package repostore

import (
	"context"

	binaryrepo "github.com/martencassel/binaryrepo"
	postgres "github.com/martencassel/binaryrepo/pkg/postgres"
)

type repoStore struct {
	db *postgres.Database
}

func NewRepoStore(db *postgres.Database) binaryrepo.RepoStore {
	rs := repoStore{db: db}
	return rs;
}

func (u repoStore) Create(ctx context.Context, repo *binaryrepo.Repo) error {
	txn, err := u.db.Begin(ctx);
	if err != nil {
		return err
	}
	err = txn.InsertRepo(*repo)
	if err != nil {
		return err
	}
	return nil
}

func (u repoStore) Get(ctx context.Context, name string) (*binaryrepo.Repo, error) {
	txn, err := u.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	repo, err := txn.GetRepo(name)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

func (u repoStore) List(ctx context.Context) ([]binaryrepo.Repo, error) {
	repos := []binaryrepo.Repo{}
	txn, err := u.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	repos, err = txn.ListRepos()
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (u repoStore) Exists(ctx context.Context, name string) (bool, error) {
	txn, err := u.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	exists, err := txn.LookupRepo(name)
	if err != nil {
		return false, err
	}
	return exists, nil
}