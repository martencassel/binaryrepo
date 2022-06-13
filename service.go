package binaryrepo

import (
	"context"
)


type UserService interface {
	Get(ctx context.Context,name string) (*User, error)
	Create(ctx context.Context,user *User) error
	List(ctx context.Context) ([]*User, error)
}