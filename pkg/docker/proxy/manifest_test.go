package dockerproxy

import (
	"os"
	"context"
	"log"
	"testing"

	client "github.com/martencassel/binaryrepo/pkg/docker/client"
	"github.com/stretchr/testify/assert"
)

// HEAD /v2/<name>/manifests/<reference>
func TestHeadManifest(t *testing.T) {
	t.Run("Check if upstream manifests exists", func(t *testing.T) {
		ctx := context.Background()
    	        hubUser := os.Getenv("DOCKERHUB_USERNAME")
                hubPass := os.Getenv("DOCKERHUB_PASSWORD")
		opt := client.Opt{
			Domain:   "docker.io",
			SkipPing: false,
			Timeout:  0,
			NonSSL:   false,
			Insecure: false,
			Debug:    false,
			Headers:  nil,
		}
		config := &client.AuthConfig{
			Username: hubUser,
			Password: hubPass,
		}
		client, err := client.New(ctx, *config, opt)
		if err != nil {
			t.Errorf("New error: %s", err)
		}
		client.SetConfig("https://registry-1.docker.io", "docker.io", config)
		err = client.Ping(ctx)
		if err != nil {
			t.Errorf("Ping error: %s", err)
		}
		exists, resp, err := client.HasManifest(ctx, "library/alpine", "latest")
		if err != nil {
			t.Errorf("HasManifest error: %s", err)
		}
		assert.True(t, exists, "found known manifest by ref as tag")
		assert.NotNil(t, resp)
		etag := resp.Header.Get("Etag")
		assert.NotNil(t, etag)
		log.Println(etag)

		exists, resp, err = client.HasManifest(ctx, "library/alpine", "sha256:686d8c9dfa6f3ccfc8230bc3178d23f84eeaf7e457f36f271ab1acc53015037c")
		etag = resp.Header.Get("Etag")
		t.Log(etag)
		assert.NotNil(t, resp)
		if err != nil {
			t.Errorf("HasManifest error: %s", err)
		}
		assert.NotNil(t, etag)
		assert.True(t, exists, "found known manifest by ref as tag ")
	})
}
