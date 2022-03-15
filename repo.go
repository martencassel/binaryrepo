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
	ID          int	`json:"id"`
	Name        string `json:"name"`
	Type        RepoType `json:"repo_type"`
	PkgType     PkgType `json:"package_type"`
	URL         string `json:"remote_url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Group       []string `json:"group"`
	Account     string `json:"account"`
	AccessToken string `json:"access_token"`
}