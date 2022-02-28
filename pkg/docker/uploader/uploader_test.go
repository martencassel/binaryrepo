package uploader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUploadManager(t *testing.T) {
	//os.RemoveAll("/tmp")
	NewUploadManager("/tmp/uploads")
	_, err := os.Stat("/tmp")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateUpload(t *testing.T) {
	u := NewUploadManager("/tmp/uploads")
	_, err := u.CreateUpload("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, u.Exists("123e4567-e89b-12d3-a456-426614174000"))
}

func TestWriteUpload(t *testing.T) {
	u := NewUploadManager("/tmp/uploads")
	_, err := u.CreateUpload("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatal(err)
	}
	u.WriteFile("123e4567-e89b-12d3-a456-426614174000", []byte("hello"))
	u.AppendFile("123e4567-e89b-12d3-a456-426614174000", []byte(" world"))
	u.AppendFile("123e4567-e89b-12d3-a456-426614174000", []byte(" world"))
	u.AppendFile("123e4567-e89b-12d3-a456-426614174000", []byte(" world"))
}
