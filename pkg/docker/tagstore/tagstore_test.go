package tagstore

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/containers/image/manifest"
)

func TestCreateTagStore(t *testing.T) {
	os.RemoveAll("/tmp/tagstore")
	NewTagStore("/tmp/tagstore")
	_, err := os.Stat("/tmp/tagstore")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateTags(t *testing.T) {
	os.RemoveAll("/tmp/tagstore")
	ts := NewTagStore("/tmp/tagstore")
	err := ts.WriteTag("ubuntu", "v1", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	// Assert that the tag exists
	if ts.TagExists("ubuntu", "v1") != true {
		t.Fatal("tag not found")
	}
	// Assert that the tag contains the correct digest
	digest, err := ts.ReadTag("ubuntu", "v1")
	if err != nil {
		t.Fatal(err)
	}
	if digest != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatal("tag not found")
	}
}

func TestGuessManifestType(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	mimeType := manifest.GuessMIMEType(b)
	if mimeType != "application/vnd.docker.distribution.manifest.v2+json" {
		t.Fatal("mime type not guessed correctly")
	}
}

func TestListImageTags(t *testing.T) {
	ts := NewTagStore("/tmp/tagstore")
	err := ts.WriteTag("ubuntu", "v1", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	err = ts.WriteTag("ubuntu", "v2", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	err = ts.WriteTag("ubuntu", "v3", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	tags, err := ts.GetTags("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 3 {
		t.Fatal("tag not found")
	}
}

func TestListImageRepos(t *testing.T) {
	ts := NewTagStore("/tmp/tagstore")
	err := ts.WriteTag("ubuntu", "v1", "13b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	err = ts.WriteTag("alpine", "v1", "21b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}
	err = ts.WriteTag("redis", "v1", "33b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	if err != nil {
		t.Fatal(err)
	}

	repos, err := ts.ListRepos()
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 3 {
		t.Fatal("repo not found")
	}

}
