package binaryrepo

type RepoType int64

const (
	Remote RepoType = 0
	Local  RepoType = 1
	Group  RepoType = 2
)

type PkgType int64

const (
	Helm   PkgType = 0
	Golang PkgType = 1
	Docker PkgType = 2
)

type Repo struct {
	ID          int
	Name        string
	Type        RepoType
	PkgType     PkgType
	URL         string
	Username    string
	Password    string
	Group       []string
	Account     string
	AccessToken string
}