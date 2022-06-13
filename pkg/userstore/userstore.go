package userstore

import (
	"context"
	"log"

	binaryrepo "github.com/martencassel/binaryrepo"
	"github.com/martencassel/binaryrepo/pkg/postgres"
)

type userStore struct {
	db *postgres.Database
}

func NewUserStore(db *postgres.Database) binaryrepo.UserStore {
	return &userStore{
		db: db,
	}
}

func (u *userStore) GetUser(ctx context.Context, username string) (*binaryrepo.User, error) {
	txn, err := u.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	user, err := txn.GetUser(username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userStore) LookupUser(ctx context.Context, username, password string) (bool, error) {
	txn, err := u.db.Begin(ctx)
	if err != nil {
		log.Println(err)
		return false, err
	}
	user := binaryrepo.User{Username: username, Password: password}
	exists, err := txn.LookupUser(user)
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = txn.Commit(ctx)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}