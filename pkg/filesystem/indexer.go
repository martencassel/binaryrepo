package filesystem

import (
	"context"
	"errors"

	"github.com/martencassel/binaryrepo"
	"github.com/opencontainers/go-digest"
)

type indexService struct {
	nodeStore binaryrepo.NodeStore
	fileStore binaryrepo.Filestore
}

func NewIndexer(ns *binaryrepo.NodeStore, fs *binaryrepo.Filestore) binaryrepo.Indexer {
	return &indexService{
		nodeStore: *ns,
		fileStore: *fs,
	}
}

func (indexer *indexService) Open(repo string, name string) (*binaryrepo.FileInfo, error) {
	node, err := indexer.nodeStore.GetNode(context.Background(), repo, name)
	if err != nil {
		return nil, err
	}
	return &binaryrepo.FileInfo{
		Repo: node.RepoName,
		Name: node.NodeName,
		Path: node.Path,
		IsFolder: node.IsFolder,
		Digest: digest.FromString(node.Checksum),
	}, nil
}


func (i *indexService) Stat(repo string, path string) (*binaryrepo.FileInfo, error) {
	// Check if the node exists, if so return file info from it.
	node, err := i.nodeStore.GetNode(context.Background(), repo, path)
	if err != nil {
		return nil, err
	}
	return &binaryrepo.FileInfo{
		Repo: node.RepoName,
		Name: node.NodeName,
		Path: node.Path,
		IsFolder: node.IsFolder,
		Digest: digest.FromString(node.Checksum),
	}, nil
}

func (i *indexService) Read(repoName string, path string, b []byte) (int, error) {
	// Check if file exists and is not a folder
	node, err := i.nodeStore.GetNode(context.Background(), repoName, path)
	if err != nil {
		return 0, err
	}
	// If it is a folder, return an error
	if node.IsFolder {
		return 0, errors.New("Cannot read a folder")
	}
	// If it is a file, read the file from the filestore using the checksum as the key
	b, err = i.fileStore.ReadFile(digest.FromString(node.Checksum))
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (i *indexService) ReadDir(repoName string, path string) ([]*binaryrepo.FileInfo, error) {
	// Check if file exists and is a folder
	// If it is a folder, read the folder from the nodestore and all direct children
	// Check if file exists and is not a folder
	node, err := i.nodeStore.GetNode(context.Background(), repoName, path)
	if err != nil {
		return nil, err
	}
	// If it is a folder, return an error
	if !node.IsFolder {
		return nil, errors.New("Cannot read a file")
	}
	return nil, nil
}

func (i *indexService) Remove(repo string, name string ) error {
	// If its a file, remove the fileinfo object
	// if its a folder, remove the folder object and all children
	return nil
}

func (i *indexService) CreateFolder(repo string, name string) *binaryrepo.FileInfo  {
	// Create a folder object
	node := &binaryrepo.Node{
		RepoName: repo,
		NodeName: name,
		Path: name,
		IsFolder: true,
	}
	err := i.nodeStore.CreateNode(context.Background(), *node)
	if err != nil {
		return nil
	}
	return &binaryrepo.FileInfo{
		Repo: node.RepoName,
		Name: node.NodeName,
		Path: node.Path,
		IsFolder: node.IsFolder,
		Digest: digest.FromString(node.Checksum),
	}
}

func (i *indexService) WriteFile(info *binaryrepo.FileInfo, bytes []byte) error {
	// Write the file to the filestore, get checksum back
	digest, err := i.fileStore.WriteFile(bytes)
	if err != nil {
		return err
	}
	// Create a fileinfo object with the checksum
	node := &binaryrepo.Node{
		RepoName: info.Repo,
		NodeName: info.Name,
		Path: info.Path,
		IsFolder: false,
		Checksum: digest.String(),
	}
	err = i.nodeStore.CreateNode(context.Background(), *node)
	if err != nil {
		return err
	}
	return nil
}