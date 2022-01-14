package client

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPingable(t *testing.T) {
	testcases := map[string]struct {
		registry Registry
		expect   bool
	}{
		"Docker": {
			registry: Registry{URL: "https://registry-1.docker.io"},
			expect:   true,
		},
		"GCR_global": {
			registry: Registry{URL: "https://gcr.io"},
			expect:   false,
		},
	}
	for label, testcatestcases := range testcases {
		if testcatestcases.registry.Pingable() != testcatestcases.expect {
			t.Fatalf("%s: expected %v, got %v", label, testcatestcases.expect, testcatestcases.registry.Pingable())
		}
	}
}

func TestPingDockerHubAuth(t *testing.T) {
	hubUser := os.Getenv("DOCKERHUB_USERNAME")
	hubPass := os.Getenv("DOCKERHUB_PASSWORD")
	if hubUser == "" || hubPass == "" {
		t.Skip("DOCKERHUB_USERNAME and DOCKERHUB_PASSWORD must be set")
	}
	ctx := context.Background()
	authConfig := &AuthConfig{
		Username:      hubUser,
		Password:      hubPass,
		Scope:         "repository:library/redis:pull",
		ServerAddress: "https://registry-1.docker.io",
	}
	opt := &Opt{
		Domain:   "docker.io",
		SkipPing: false,
		Timeout:  time.Second * 120,
		NonSSL:   false,
		Headers:  map[string]string{},
	}
	registry, err := New(ctx, *authConfig, *opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := registry.Ping(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPingACRAuth(t *testing.T) {
	acrUser := os.Getenv("ACR_USERNAME")
	acrPass := os.Getenv("ACR_PASSWORD")
	acrServerAddress := os.Getenv("ACR_SERVER")
	if acrUser == "" || acrPass == "" || acrServerAddress == "" {
		t.Skip("ACR_USERNAME, ACR_PASSWORD and ACR_SERVER must be set")
	}
	ctx := context.Background()
	authConfig := &AuthConfig{
		Username:      acrUser,
		Password:      acrPass,
		ServerAddress: fmt.Sprintf("https://%s", acrServerAddress),
		Scope:         "repository:redis:push,pull",
	}
	opt := &Opt{
		Domain:   acrServerAddress,
		SkipPing: false,
		Timeout:  time.Second * 120,
		NonSSL:   false,
		Headers:  map[string]string{},
	}
	registry, err := New(ctx, *authConfig, *opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := registry.Ping(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
