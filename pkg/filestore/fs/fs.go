package filestore

import (
	"fmt"
	"io/ioutil"
	"os"

	binaryrepo "github.com/martencassel/binaryrepo"
	digest "github.com/opencontainers/go-digest"
	log "github.com/rs/zerolog/log"
)

//var _ FileStore = &fileStore{}

type fileStore struct {
	BasePath string
}

func NewFileStore(basePath string) binaryrepo.Filestore {
	fs := fileStore{
		BasePath: basePath,
	}
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err = os.MkdirAll(basePath+"/uploads", 0755)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return &fs
}

func (fs *fileStore) GetBasePath() string {
	return fs.BasePath
}

func getFilePath(basePath string, digest digest.Digest) (string, string, string) {
	encodedDigest := digest.Encoded()
	first := encodedDigest[0:2]
	last := encodedDigest[2:]
	folderName := string(first)
	fileName := string(last)
	filePath := fmt.Sprintf("%s/%s/%s", basePath, folderName, fileName)
	folderPath := fmt.Sprintf("%s/%s", basePath, folderName)
	return filePath, folderPath, fileName
}

func (fs *fileStore) Remove(digest digest.Digest) error {
	if !fs.Exists(digest) {
		return nil
	}
	filePath, _, _ := getFilePath(fs.BasePath, digest)
	err := os.Remove(filePath)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return err
	}
	return nil
}

func (fs *fileStore) Exists(digest digest.Digest) bool {
	filePath, _, _ := getFilePath(fs.BasePath, digest)
	_, err := os.Stat(filePath)
	return err == nil
}

func (fs *fileStore) WriteFile(b []byte) (digest.Digest, error) {
	digest := digest.FromBytes(b)
	filePath, folderPath, _ := getFilePath(fs.BasePath, digest)
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return digest.Algorithm().FromString(""), err
	}
	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return digest.Algorithm().FromString(""), err
	}
	return digest, nil
}

func (fs *fileStore) ReadFile(digest digest.Digest) ([]byte, error) {
	if !fs.Exists(digest) {
		return nil, nil
	}
	filePath, _, _ := getFilePath(fs.BasePath, digest)
	file, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	return b, nil
}