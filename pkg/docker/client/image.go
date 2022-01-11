package client

import (
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
)

// Image holds information about an image.
type Image struct {
	Domain string
	Path   string
	Tag    string
	Digest digest.Digest
	named  reference.Named
}

// String returns the string representation of an image.
func (i Image) String() string {
	return i.named.String()
}

// Reference returns either the digest if it is non-empty or the tag for the image.
func (i Image) Reference() string {
	if len(i.Digest.String()) > 1 {
		return i.Digest.String()
	}
	return i.Tag
}

// WithDigest sets the digest for an image.
func (i *Image) WithDigest(digest digest.Digest) (err error) {
	i.Digest = digest
	i.named, err = reference.WithDigest(i.named, digest)
	return err
}
