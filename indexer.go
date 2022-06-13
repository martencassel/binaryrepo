package binaryrepo

import (
	"github.com/opencontainers/go-digest"
)

type FileInfo struct {
	Name string
	Path string
	Repo string
	IsFolder bool
	Digest digest.Digest
	Children []*FileInfo
}

func (f *FileInfo) IsDir() bool {
	return f.IsFolder
}

func (f *FileInfo) IsFile() bool {
	return !f.IsFolder
}

// Indexer is an interface that can process file metadata operations.
type Indexer interface {
	Open(repo string, name string) (*FileInfo, error)
	Stat(repo string, path string) (*FileInfo, error)
	Read(repo string, path string, b []byte) (int, error)
	ReadDir(repoName string, path string) ([]*FileInfo, error)
	Remove(repo string, name string ) error
	WriteFile(info *FileInfo, bytes []byte) error
	CreateFolder(repo string, name string) *FileInfo
}