package filestore

import (
	"io/ioutil"
	"os"
	"testing"

	_ "crypto/sha256"

	"github.com/opencontainers/go-digest"
)

func TestFileStore(t *testing.T) {
	t.Run("Get filepath from cheksum", func(t *testing.T) {
		os.RemoveAll("/tmp/filestore")
		fs := NewFileStore("/tmp/filestore")
		b, _ := ioutil.ReadFile("./testdata/file1")
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		path := getFilePath(fs.BasePath, digest)
		if path != "/tmp/filestore/c1/47efcfc2d7ea666a9e4f5187b115c90903f0fc896a56df9a6ef5d8f3fc9f31" {
			t.Errorf("got %s, want %s", path, "/tmp/filestore/c1/47efcfc2d7ea666a9e4f5187b115c90903f0fc896a56df9a6ef5d8f3fc9f31")
		}
	})

	t.Run("Create filestore at base directory", func(t *testing.T) {
		os.RemoveAll("/tmp/filestore")
		fs := NewFileStore("/tmp/filestore")
		_, err := os.Stat(fs.BasePath)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Writing a file", func(t *testing.T) {
		os.RemoveAll("/tmp/filestore")
		fs := NewFileStore("/tmp/filestore")
		b, _ := ioutil.ReadFile("./testdata/file1")
		fsDigest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		digest := digest.FromBytes(b)
		if fsDigest != digest {
			t.Fatal("sha256 mismatch")
		}
		out, err := fs.ReadFile(digest)
		if err != nil {
			t.Fatal(err)
		}
		if string(out) != string(b) {
			t.Fatal("read mismatch")
		}
		if fs.Exists(digest) != true {
			t.Fatal("file not found")
		}
	})
	t.Run("Check for existance of file using digest", func(t *testing.T) {
		os.RemoveAll("/tmp/filestore")
		fs := NewFileStore("/tmp/filestore")
		sha := "sha256:56704d8d370580ad16fcfbf725982551da20fb82b4450f9aedfd055fa9857967"
		exists := fs.Exists(digest.Digest(sha))
		if exists == true {
			t.Fatal("file found")
		}
	})
	t.Run("Read a non existing file", func(t *testing.T) {
		os.RemoveAll("/tmp/filestore")
		fs := NewFileStore("/tmp/filestore")
		sha := "sha256:56704d8d370580ad16fcfbf725982551da20fb82b4450f9aedfd055fa9857967"
		b, err := fs.ReadFile(digest.Digest(sha))
		if b != nil && err == nil {
			t.Fatal("Not expecting file to be found")
		}
	})
}
