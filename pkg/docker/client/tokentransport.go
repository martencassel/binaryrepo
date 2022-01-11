package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type TokenTransport struct {
	Transport http.RoundTripper
	Username  string
	Password  string
	Account   string
	Scope     string
}

func (t *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	authService, err := isTokenDemand(resp)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	if authService == nil {
		return resp, nil
	}
	resp.Body.Close()
	return t.authAndRetry(authService, req)
}

type authService struct {
	Realm   *url.URL
	Account string
	Service string
	Scope   []string
}

func (a *authService) Request(username, password, scope string) (*http.Request, error) {
	q := a.Realm.Query()
	q.Set("service", a.Service)
	q.Set("account", username)
	if scope != "" {
		q.Set("scope", scope)
	} else {
		for _, s := range a.Scope {
			q.Set("scope", s)
		}
	}
	a.Realm.RawQuery = q.Encode()
	req, err := http.NewRequest("GET", a.Realm.String(), nil)
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}
	return req, err
}

func isTokenDemand(resp *http.Response) (*authService, error) {
	if resp.StatusCode != http.StatusUnauthorized {
		return nil, nil
	}
	return parseAuthHeader(resp.Header)
}

type authToken struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
}

func (t authToken) String() (string, error) {
	if t.Token != "" {
		return t.Token, nil
	}
	if t.AccessToken != "" {
		return t.AccessToken, nil
	}
	return "", errors.New("auth token cannot be empty")
}

func (t *TokenTransport) authAndRetry(authService *authService, req *http.Request) (*http.Response, error) {
	token, authResp, err := t.auth(req.Context(), authService)
	if err != nil {
		return authResp, err
	}
	response, err := t.retry(req, token)
	if response != nil {
		response.Header.Set("request-token", token)
	}
	return response, err
}

func (t *TokenTransport) auth(ctx context.Context, authService *authService) (string, *http.Response, error) {
	authReq, err := authService.Request(t.Username, t.Password, t.Scope)
	if err != nil {
		return "", nil, err
	}
	c := http.Client{
		Transport: t.Transport,
	}
	resp, err := c.Do(authReq.WithContext(ctx))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", resp, err
	}
	var authToken authToken
	if err := json.NewDecoder(resp.Body).Decode(&authToken); err != nil {
		return "", nil, err
	}
	token, err := authToken.String()
	return token, nil, err
}

func (t *TokenTransport) retry(req *http.Request, token string) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return t.Transport.RoundTrip(req)
}
