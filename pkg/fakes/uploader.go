package fakes

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/martencassel/binaryrepo"
)

type uploader struct {
	indexer *binaryrepo.Indexer
}

func NewUploader(idx binaryrepo.Indexer) binaryrepo.Uploader {
	return &uploader{
		indexer: &idx,
	}
}

func (u *uploader)   CreateUpload() (uuid.UUID, error) {
	// Create a directory
	uuid, err := uuid.NewUUID()

	_ = os.MkdirAll("/tmp/uploads", 0755)
	f,_  := os.Create(fmt.Sprintf("/tmp/uploads/%s", uuid.String()))
	defer f.Close()
	return uuid, err
}

func (u *uploader) ReadUpload(uuid string) ([]byte, error) {
	// Read the file
	f, err := os.Open("/tmp/uploads/" + uuid)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bytes := make([]byte, 1024)
	n, err := f.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes[:n], nil
}

func (u *uploader) WriteFile(uuid string, bytes []byte) error {
	f, err := os.Create("/tmp/uploads/" + uuid)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func (u *uploader) Exists(uuid string) (bool) {
	exists, _ := fileExists("/tmp/uploads/" + uuid)
	return exists
}

func (u *uploader) Remove(uuid string) error {
	// Remove upload file
	return os.Remove("/tmp/uploads/" + uuid)
}

func (u *uploader) AppendFile(uuid string, bytes []byte) error {
	return nil
}

func fileExists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
}