package memory

import (
	_ "crypto/sha256"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	t.Run("Writing a file", func(t *testing.T) {
		fs := NewFileStore("/tmp/memory")
		b := []byte("hello world")
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		if !fs.Exists(digest) {
			t.Fatal("file not found")
		}
	})
	t.Run("Remove a file", func(t *testing.T) {
		fs := NewFileStore("/tmp/memory")
		b := []byte("hello world")
		digest, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		err = fs.Remove(digest)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fs.Exists(digest), false)
	})
	t.Run("Check for existance of file using digest", func(t *testing.T) {
		fs := NewFileStore("/tmp/memory")
		b := []byte("hello world")
		_, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, fs.Exists("sha256:b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"), true)
	})
	t.Run("Read a non existing file", func(t *testing.T) {
		fs := NewFileStore("/tmp/memory")
		b := []byte("hello world")
		_, err := fs.WriteFile(b)
		if err != nil {
			t.Fatal(err)
		}
		sha := "sha256:56704d8d370580ad16fcfbf725982551da20fb82b4450f9aedfd055fa9857967"
		err = fs.Remove(digest.FromString(sha))
		if err != nil {
			t.Fatal(err)
		}
		b, err = fs.ReadFile(digest.Digest(sha))
		if b != nil && err == nil {
			t.Fatal("Not expecting file to be found")
		}
		assert.Equal(t, fs.Exists("sha256:56704d8d370580ad16fcfbf725982551da20fb82b4450f9aedfd055fa9857967"), false)
	})
}
