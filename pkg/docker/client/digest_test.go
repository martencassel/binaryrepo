package client

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestDigestFromDockerHub(t *testing.T) {
	hubUser := os.Getenv("DOCKERHUB_USERNAME")
	hubPass := os.Getenv("DOCKERHUB_PASSWORD")
	if hubUser == "" || hubPass == "" {
		t.Skip("DOCKERHUB_USERNAME and DOCKERHUB_PASSWORD must be set")
	}
	ctx := context.Background()
	r, err := New(ctx, AuthConfig{
		Username:      hubUser,
		Password:      hubPass,
		Scope:         "repository:library/redis:pull",
		ServerAddress: "https://registry-1.docker.io",
	}, Opt{
		Domain:   "docker.io",
		SkipPing: false,
		Timeout:  time.Second * 120,
		NonSSL:   false,
		Insecure: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d, _, err := r.Digest(ctx, Image{Domain: "docker.io", Path: "library/redis", Tag: "latest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == "" {
		t.Error("empty digest received")
	}
}
