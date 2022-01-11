package repo

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

var _repos []Repo = []Repo{
	{
		ID:      1,
		Name:    "golang-remote",
		Type:    Remote,
		PkgType: Golang,
		URL:     "https://proxy.golang.org",
		auth:    false,
	},
	{
		ID:      2,
		Name:    "helm-bitnami-remote",
		Type:    Remote,
		PkgType: Helm,
		URL:     "https://charts.bitnami.com/bitnami",
		auth:    false,
	},
	{
		ID:       3,
		Name:     "docker-remote",
		Type:     Remote,
		PkgType:  Docker,
		URL:      "https://registry-1.docker.io",
		auth:     true,
		Username: "",
		Password: "",
		Account:  "",
	},
	{
		ID:       4,
		Name:     "docker-remote2",
		Type:     Remote,
		PkgType:  Docker,
		URL:      "https://registry-1.docker.io",
		auth:     true,
		Username: "",
		Password: "",
	},
	{
		ID:      5,
		Name:    "docker-local",
		Type:    Local,
		PkgType: Docker,
		auth:    false,
	},
	{
		ID:      6,
		Name:    "helm-local",
		Type:    Local,
		PkgType: Helm,
		auth:    false,
	},
	{
		ID:      7,
		Name:    "helm-group",
		Type:    Group,
		PkgType: Helm,
		Group: []string{
			"helm-local",
			"helm-bitnami-remote",
		},
	},
}

type RepoIndex struct {
	Repos []Repo
}

type Repo struct {
	ID          int
	Name        string
	Type        RepoType
	PkgType     PkgType
	URL         string
	Username    string
	Password    string
	auth        bool
	Group       []string
	Account     string
	AccessToken string
}

func NewRepoIndex() *RepoIndex {
	return &RepoIndex{}
}

func (index *RepoIndex) GetRepos() []Repo {
	return index.Repos
}

func (index *RepoIndex) FindRepo(name string) *Repo {
	for i, repo := range index.Repos {
		if repo.Name == name {
			return &index.Repos[i]
		}
	}
	return nil
}

func (index *RepoIndex) AddRepo(repo Repo) {
	index.Repos = append(index.Repos, repo)
}
