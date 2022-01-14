package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/opencontainers/go-digest"
)

func TestLayerFromDockerHub(t *testing.T) {
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
	exists, _, err := r.HasLayer(ctx, "library/redis", "sha256:c8388a79482fce47e8f9cc1811df4f4fbd12260fee9128b29903bf4a3f33dd01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("layer does not exist")
	}
	digest, err := digest.Parse("sha256:c8388a79482fce47e8f9cc1811df4f4fbd12260fee9128b29903bf4a3f33dd01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rc, _, err := r.DownloadLayer(ctx, "library/redis", digest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc == nil {
		t.Error("empty layer received")
	}
}
