package client

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	t.Run("Test client", func(t *testing.T) {
		ctx := context.Background()
		hubUser := os.Getenv("DOCKERHUB_USERNAME")
                hubPass := os.Getenv("DOCKERHUB_PASSWORD")
		opt := Opt{
			Domain:   "docker.io",
			SkipPing: false,
			Timeout:  0,
			NonSSL:   false,
			Insecure: false,
			Debug:    false,
			Headers:  nil,
		}
		config := &AuthConfig{
			Username: hubUser,
			Password: hubPass,
		}
		client, err := New(ctx, *config, opt)
		if err != nil {
			t.Errorf("New error: %s", err)
		}
		client.SetConfig("https://registry-1.docker.io", "docker.io", config)
		err = client.Ping(ctx)
		if err != nil {
			t.Errorf("Ping error: %s", err)
		}
		_, resp, err := client.HasManifest(ctx, "library/alpine", "latest")
		if err != nil {
			t.Errorf("DownloadLayer error: %s", err)
		}
		out, err := os.Create("/tmp/alpine.tar")
		if err != nil {
			t.Errorf("Create error: %s", err)
		}
		io.Copy(out, resp.Body)

	// 	regclient := NewRegistryClient()

	// 	r, err := regclient.New(ctx, regclient.AuthConfig{
	// 		Username:      username,
	// 		Password:      password,
	// 		Scope:         scope,
	// 		ServerAddress: srvaddr,
	// }, regclient.Opt{
	// 		Domain:   domain,
	// 		SkipPing: false,
	// 		Timeout:  time.Minute * 10,
	// 		NonSSL:   false,
	// 		Insecure: false,
	// })

	})
}
