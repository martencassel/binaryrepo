package client

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestDigestFromDockerHub(t *testing.T) {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	hubUser := os.Getenv("DOCKERHUB_USERNAME")
	hubPass := os.Getenv("DOCKERHUB_PASSWORD")
	ctx := context.Background()
	r, err := New(ctx, AuthConfig{
		Username:      hubUser,
		Password:      hubPass,
		Scope:         "repository:library/redis:pull",
		ServerAddress: "https://registry-1.docker.io",
	}, Opt{
		Domain:   "docker.io",
		SkipPing: false,
		Timeout:  time.Second * 30,
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
