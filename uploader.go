package binaryrepo

// Uploader is an interface for a service that supports docker registry v2 upload operations.
type Uploader interface {
	CreateUpload(uuid string) (string, error)
	ReadUpload(uuid string) ([]byte, error)
	WriteFile(uuid string, bytes []byte) error
	AppendFile(uuid string, bytes []byte) error
	Exists(uuid string) bool
	Remove(uuid string) error
}