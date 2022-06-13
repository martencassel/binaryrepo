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
	ID 			   int	  `json:"id"`
	Name 	 	   string `json:"name"`
	Repotype 	   string `json:"repo_type"`
	Pkgtype   	   string `json:"pkg_type"`
	Remoteurl      string `json:"remote_url"`
	RemoteUsername string `json:"remote_username"`
	RemotePassword string `json:"remote_password"`
	Anonymous      bool `json:"anonymous"`
}
