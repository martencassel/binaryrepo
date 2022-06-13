package fakes

import (
	"github.com/martencassel/binaryrepo"
	"github.com/opencontainers/go-digest"
)


type TagKey struct {
	repo string
	path string
}

type TagInfo struct {
	name string
	digest digest.Digest
}

type tagstore struct {
	tags map[TagKey]*TagInfo
}

func NewTagStore() binaryrepo.TagStore {
	return &tagstore{
		tags: make(map[TagKey]*TagInfo, 1000),
	}
}

func (t *tagstore) Exists(repo string, tag string) bool {
	return t.tags[TagKey{repo, tag}] != nil
}

func (t *tagstore) GetTag(repo, tag string) (digest.Digest, error) {
	return t.tags[TagKey{repo, tag}].digest, nil
}

func (t *tagstore) GetTags(repo string) ([]string, error) {
	tags := make([]string, 0)
	for key, _ := range t.tags {
		if key.repo == repo {
			tags = append(tags, key.path)
		}
	}
	return tags, nil
}

func (t *tagstore) WriteTag(repo string, tag string, digest digest.Digest) error {
	t.tags[TagKey{repo, tag}] = &TagInfo{
		name: tag,
		digest: digest,
	}
	return nil
}