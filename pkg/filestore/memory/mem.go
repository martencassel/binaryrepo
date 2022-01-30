package memory

import (
	"bytes"
	"errors"

	_ "crypto/sha256"

	"github.com/opencontainers/go-digest"
)

type FileStore struct {
	BasePath string
	fileMap  map[digest.Digest]*bytes.Buffer
}

func NewFileStore(basePath string) *FileStore {
	fs := &FileStore{
		BasePath: basePath,
		fileMap:  make(map[digest.Digest]*bytes.Buffer),
	}
	return fs
}

func (fs *FileStore) ReadFile(digest digest.Digest) ([]byte, error) {
	if ok := fs.fileMap[digest]; ok != nil {
		return []byte(digest.Hex()), nil
	}
	return nil, errors.New("file not found")
}

func (fs *FileStore) WriteFile(b []byte) (digest.Digest, error) {
	fs.fileMap[digest.FromBytes(b)] = bytes.NewBuffer(b)
	return digest.FromBytes(b), nil
}

func (fs *FileStore) Remove(digest digest.Digest) error {
	delete(fs.fileMap, digest)
	return nil
}

func (fs *FileStore) Exists(digest digest.Digest) bool {
	_, ok := fs.fileMap[digest]
	return ok
}
