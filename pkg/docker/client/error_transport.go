package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpStatusError struct {
	Response *http.Response
	Body     []byte
}

func (err *httpStatusError) Error() string {
	return fmt.Sprintf("http: non-successful response (status=%v body=%q)", err.Response.StatusCode, err.Body)
}

var _ error = &httpStatusError{}

type ErrorTransport struct {
	Transport http.RoundTripper
}

func (t *ErrorTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(request)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 500 || resp.StatusCode == http.StatusUnauthorized {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http: failed to read response body (status=%v, err=%q)", resp.StatusCode, err)
		}
		return nil, &httpStatusError{
			Response: resp,
			Body:     body,
		}
	}
	return resp, err
}
