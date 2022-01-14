package repo

import (
	log "github.com/rs/zerolog/log"
)

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
	Group       []string
	Account     string
	AccessToken string
}

func NewRepoIndex() *RepoIndex {
	return &RepoIndex{}
}

func (index *RepoIndex) AddRepo(repo Repo) {
	index.Repos = append(index.Repos, repo)
}

func (index *RepoIndex) GetRepos() []Repo {
	log.Info().Msgf("%v", index)
	return index.Repos
}

func (index *RepoIndex) FindRepo(name string) *Repo {
	log.Info().Msgf("FindRepo: ", name)
	log.Info().Msgf("%v", index.Repos)
	for i, repo := range index.Repos {
		if repo.Name == name {
			return &index.Repos[i]
		}
	}
	return nil
}
