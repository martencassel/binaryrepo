package binaryrepo

import "context"

type UserStore interface {
	GetUser(ctx context.Context, username string) (*User, error)
	LookupUser(ctx context.Context, username, password string) (bool, error)
}