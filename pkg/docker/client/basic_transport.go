package client

import (
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
	//////log.Info().Msgf("%s %s", req.URL.String(), req.Header.Get("Authorization"))
	if strings.HasPrefix(req.URL.String(), t.URL) && req.Header.Get("Authorization") == "" {
		if t.Username != "" || t.Password != "" {
			req.SetBasicAuth(t.Username, t.Password)
		}
	}
	resp, err := t.Transport.RoundTrip(req)
	// //log.Info().Msg("Basic transport:")
	// ////log.Info().Msgf("%v", resp)
	// ////log.Info().Msgf("%v", resp.StatusCode)
	return resp, err
}
