package dockerproxy

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	authUrlDockerHub string = "https://auth.docker.io/token?account=%s&client_id=docker&offline_token=true&scope=repository:%s:pull&service=registry.docker.io"
)

func GetAuthResponse(account, username, password, scope string) (accessToken string, err error) {
	//	log.Printf("GetAuthResponse(%s, %s, %s, %s)", account, username, password, scope)
	url := fmt.Sprintf(authUrlDockerHub, account, scope)
	//	log.Printf("GetAuthResponse: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error authenticating to %s: %s", url, resp.Status)
	}
	var token struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}
	return token.Token, nil
}
