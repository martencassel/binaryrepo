package client

import (
	"log"
	"net/http"
	"strings"
)

type BasicTransport struct {
	Transport http.RoundTripper
	URL       string
	Username  string
	Password  string
}

func (t *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Printf("%s %s", req.URL.String(), req.Header.Get("Authorization"))
	if strings.HasPrefix(req.URL.String(), t.URL) && req.Header.Get("Authorization") == "" {
		if t.Username != "" || t.Password != "" {
			req.SetBasicAuth(t.Username, t.Password)
		}
	}
	resp, err := t.Transport.RoundTrip(req)
	return resp, err
}
