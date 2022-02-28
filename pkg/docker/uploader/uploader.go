package uploader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type UploadManager struct {
	BasePath string
}

func NewUploadManager(basePath string) *UploadManager {
	ts := &UploadManager{
		BasePath: basePath,
	}
	err := os.MkdirAll(ts.BasePath, 0755)
	if err != nil {
		panic(err)
	}
	return ts
}

func (u *UploadManager) CreateUpload(uuid string) (string, error) {
	filename := filepath.Join(u.BasePath, uuid)
	err := ioutil.WriteFile(filename, []byte{}, 0644)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create upload %s", uuid)
	}
	return filename, nil
}

func (u *UploadManager) ReadUpload(uuid string) ([]byte, error) {
	filename := filepath.Join(u.BasePath, uuid)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read upload %s", uuid)
	}
	return b, nil
}

func (u *UploadManager) WriteFile(uuid string, bytes []byte) error {
	filename := filepath.Join(u.BasePath, uuid)
	err := ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return errors.Errorf("Failed to write file %s", filename)
	}
	return nil
}

func (u *UploadManager) AppendFile(uuid string, bytes []byte) error {
	filename := filepath.Join(u.BasePath, uuid)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	if _, err = f.Write(bytes); err != nil {
		panic(err)
	}
	return nil
}

func (u *UploadManager) Exists(uuid string) bool {
	filePath := filepath.Join(u.BasePath, uuid)
	_, err := os.Stat(filePath)
	return err == nil
}

// Status of upload.
//
// Input: <uuid>
// Impl:
//		1. Open the <uuid> file.
//		2. Read the file size
// 		3. Return the offset = file size - 1
// Output: Range: bytes=0-<offset>

// Monolithic upload.
//
// Write single chunk to upload file.
// Input: <uuid>, <digest>, <Layer Binary Data>, <Size of Layer Binary Data>
// Output: The layer is written to the upload file <uuid>

// Chunked upload.
//
// Write chunk to upload file using range header offset.
// Input: <uuid>, <content-length>, <layer binary data>
