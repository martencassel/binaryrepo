package repo

import (
	"testing"
)

func TestRepoInit(t *testing.T) {
	index := NewRepoIndex()
	remote := index.FindRepo("docker-remote")
	if remote == nil {
		t.Fatal("Repo should exist")
	}
}
