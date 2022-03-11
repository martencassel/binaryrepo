package binaryrepo

// Uploader is an interface for a service that supports docker registry v2 upload operations.
type Uploader interface {

	// Create a new uploads and return the upload uuid.
	CreateUpload(uuid string) (string, error)

	// Read the contents of the upload.
	ReadUpload(uuid string) ([]byte, error)

	// Write a file to a upload.
	WriteFile(uuid string, bytes []byte) error

	// Append bytes to an existing upload.
	AppendFile(uuid string, bytes []byte) error

	// Finish the upload.
	Exists(uuid string) bool

	// Remove an upload.
	Remove(uuid string) error
}