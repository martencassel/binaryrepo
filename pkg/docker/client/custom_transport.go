package client

import "net/http"

type CustomTransport struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

func (t *CustomTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if len(t.Headers) != 0 {
		for header, value := range t.Headers {
			request.Header.Add(header, value)
		}
	}
	resp, err := t.Transport.RoundTrip(request)
	return resp, err
}
