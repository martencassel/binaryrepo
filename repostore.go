package binaryrepo

import "context"

type RepoStore interface {
	Create(ctx context.Context, repo *Repo) error
	Get(ctx context.Context, name string) (*Repo, error)
	List(ctx context.Context) ([]Repo, error)
	Exists(ctx context.Context, name string) (bool, error)
}