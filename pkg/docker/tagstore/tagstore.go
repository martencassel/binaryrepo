package tagstore

import (
	"fmt"
	"io/ioutil"
	"os"
)

type TagStore struct {
	BasePath string
}

func NewTagStore(basePath string) *TagStore {
	ts := &TagStore{
		BasePath: basePath,
	}
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		panic(err)
	}
	return ts
}

// Write a tag with the given name and reference that contains a string foo
func (ts *TagStore) WriteTag(name, reference, digest string) error {
	// Ensure that the tag directory exists
	err := os.MkdirAll(fmt.Sprintf("%s/%s", ts.BasePath, name), 0755)
	if err != nil {
		panic(err)
	}
	return ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", ts.BasePath, name, reference), []byte(digest), 0644)
}

// Read tag from the tag store
func (ts *TagStore) ReadTag(name, reference string) (string, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", ts.BasePath, name, reference))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (ts *TagStore) TagExists(name, reference string) bool {
	path := fmt.Sprintf("%s/%s/%s", ts.BasePath, name, reference)
	_, err := os.Stat(path)
	return err == nil
}

func (ts *TagStore) GetTags(name string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", ts.BasePath, name)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, f := range files {
		tags = append(tags, f.Name())
	}
	return tags, nil
}

func (ts *TagStore) ListRepos() ([]string, error) {
	files, err := ioutil.ReadDir(ts.BasePath)
	if err != nil {
		return nil, err
	}
	var repos []string
	for _, f := range files {
		repos = append(repos, f.Name())
	}
	return repos, nil
}
