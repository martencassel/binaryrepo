package filestore

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	digest "github.com/opencontainers/go-digest"
)

type FileStore struct {
	BasePath string
}

func NewFileStore(basePath string) *FileStore {
	fs := &FileStore{
		BasePath: basePath,
	}
	os.MkdirAll(basePath, 0755)
	os.MkdirAll(basePath+"/uploads", 0755)
	return fs
}

func getFilePath(basePath string, digest digest.Digest) string {
	encodedDigest := digest.Encoded()
	first := encodedDigest[0:2]
	last := encodedDigest[2:]
	folderName := string(first)
	fileName := string(last)
	filePath := fmt.Sprintf("%s/%s/%s", basePath, folderName, fileName)
	return filePath
}

func (fs *FileStore) Remove(digest digest.Digest) error {
	if !fs.Exists(digest) {
		return nil
	}
	encodedDigest := digest.Encoded()
	first := encodedDigest[0:2]
	last := encodedDigest[2:]
	folderName := string(first)
	fileName := string(last)
	filePath := fmt.Sprintf("%s/%s/%s", fs.BasePath, folderName, fileName)
	err := os.Remove(filePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (fs *FileStore) Exists(digest digest.Digest) bool {
	filePath := getFilePath(fs.BasePath, digest)
	_, err := os.Stat(filePath)
	return err == nil
}

func (fs *FileStore) WriteFile(b []byte) (digest.Digest, error) {
	digest := digest.FromBytes(b)
	encodedDigest := digest.Encoded()
	first := encodedDigest[0:2]
	last := encodedDigest[2:]
	folderName := string(first)
	fileName := string(last)
	folderPath := fmt.Sprintf("%s/%s", fs.BasePath, folderName)
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Fatal(err)
		return digest.Algorithm().FromString(""), err
	}
	filePath := fmt.Sprintf("%s/%s", folderPath, fileName)
	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Fatal(err)
		return digest.Algorithm().FromString(""), err
	}
	return digest, nil
}

func (fs *FileStore) ReadFile(digest digest.Digest) ([]byte, error) {
	encodedDigest := digest.Encoded()
	first := encodedDigest[0:2]
	last := encodedDigest[2:]
	folderName := string(first)
	fileName := string(last)
	filePath := fmt.Sprintf("%s/%s/%s", fs.BasePath, folderName, fileName)
	// Check if file exists
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return b, nil
}

// func (fs *FileStore) Size(digest digest.Digest) (int, error) {

// }
