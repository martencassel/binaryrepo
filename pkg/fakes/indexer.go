package fakes

import (
	"path"

	"github.com/martencassel/binaryrepo"
)

type FileKey struct {
	repo string
	path string
}

type indexer struct {
	files map[FileKey]*binaryrepo.FileInfo
}

func NewIndexer() binaryrepo.Indexer {
	return &indexer{
		files: make(map[FileKey]*binaryrepo.FileInfo),
	}
}

func (indexer *indexer) Open(repo string, name string) (*binaryrepo.FileInfo, error) {
	return indexer.files[FileKey{repo, name}], nil
}

// Create a folder
func (indexer *indexer) CreateFolder(repo string, name string) *binaryrepo.FileInfo {
	folder := indexer.files[FileKey{repo, name}]
	if folder == nil {
		folder = &binaryrepo.FileInfo{
			Repo: repo,
			Name: name,
			Path: name,
			IsFolder: true,
		}
	}
	return folder
}

// Read a file
func (indexer *indexer) Read(repo string, path string, b []byte) (int, error) {
	return 0, nil
}

func (indexer *indexer) getFolder(repo string, dirPath string) (*binaryrepo.FileInfo, error) {
	folder := indexer.CreateFolder(repo, dirPath)
	return folder, nil
}

// Read a directory
func (indexer *indexer) ReadDir(repo string, path string) ([]*binaryrepo.FileInfo, error) {
	return nil, nil
}

// Remove a file
func (indexer *indexer) Remove(repo string, path string) error {
	return nil
}

// Stat a file
func (indexer *indexer) Stat(repo string, path string) (*binaryrepo.FileInfo, error) {
	return nil, nil
}

// Write a file
func (indexer *indexer) WriteFile(info *binaryrepo.FileInfo, b []byte) error {
	folder, _ := indexer.getFolder(info.Repo, path.Dir(info.Path))
	fi := &binaryrepo.FileInfo{
		Repo: info.Repo,
		Name: info.Name,
		Path: info.Path,
		IsFolder: false,
		Digest: info.Digest,
	}
	folder.Children = append(folder.Children, fi)
	return nil
}