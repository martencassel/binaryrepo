package api

import (
	"github.com/martencassel/binaryrepo"
)

type ApiHandler struct {
	rs binaryrepo.RepoStore
	us binaryrepo.UserStore
	fs 		binaryrepo.Filestore
}

func NewApiHandler(rs binaryrepo.RepoStore,
				   us binaryrepo.UserStore,
				   fs binaryrepo.Filestore) * ApiHandler {
	return &ApiHandler {
		rs: rs,
		fs: fs,
	}
}
