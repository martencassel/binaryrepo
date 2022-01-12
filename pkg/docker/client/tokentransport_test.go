package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrBasicAuth(t *testing.T) {
	ctx := context.Background()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("www-authenticate", `Basic realm="Registry Realm",service="Docker registry"`)
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer ts.Close()
	authConfig := AuthConfig{
		Username:      "uname",
		Password:      "pword",
		ServerAddress: ts.URL,
	}
	r, err := New(ctx, authConfig, Opt{Insecure: true, Debug: true})
	if err != nil {
		t.Fatalf("expected no error creating client, got %v", err)
	}
	token, err := r.Token(ctx, ts.URL)
	if err != ErrBasicAuth {
		t.Fatalf("expected ErrBasicAuth, got %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty token, got %q", token)
	}
}
