package uploader

import (
	"github.com/google/uuid"

	"github.com/martencassel/binaryrepo"
)

type uploaderType struct {
	indexer *binaryrepo.Indexer
}

func NewUploader(idx binaryrepo.Indexer) binaryrepo.Uploader {
	return &uploaderType{
		indexer: &idx,
	}
}

func (u *uploaderType)  CreateUpload() (uuid.UUID, error)  {
	uuid, err := uuid.NewUUID()
	return uuid, err
}

func (u *uploaderType) ReadUpload(uuid string) ([]byte, error) {
	return nil, nil
}

func (u *uploaderType) WriteFile(uuid string, bytes []byte) error {
	return nil
}

func (u *uploaderType) Exists(uuid string) bool {
	return false
}

func (u *uploaderType) Remove(uuid string) error {
	return nil
}

func (u *uploaderType) AppendFile(uuid string, bytes []byte) error {
	return nil
}


