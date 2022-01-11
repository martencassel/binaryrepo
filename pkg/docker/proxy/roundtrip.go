package dockerproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/golang-jwt/jwt"
	repo "github.com/martencassel/binaryrepo/pkg/repo"
)

type Transport struct {
	username    string
	password    string
	account     string
	scope       string
	accessToken string
	next        http.RoundTripper
}

func NewReAuthTransport(repo *repo.Repo, scope, accessToken string) *Transport {
	return &Transport{
		username:    repo.Username,
		password:    repo.Password,
		account:     repo.Account,
		scope:       scope,
		accessToken: accessToken,
		next:        http.DefaultTransport,
	}
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, _err error) {
	log.Printf("RoundTrip: %s %s", req.Method, req.URL.Path)
	_, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return nil, err
	}
	for count := 0; count < 3; count++ {
		if t.accessToken != "" && checkToken(t.accessToken) {
			log.Printf("Already found access token!!!")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.accessToken))
		}
		resp, err := t.next.RoundTrip(req)
		log.Printf("%s %s %s", resp.Request.Method, resp.Request.URL.Path, resp.Status)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		_, err = httputil.DumpResponse(resp, false)
		if err != nil {
			return nil, err
		}
		log.Print("Response:", resp.StatusCode)
		if strings.Contains(req.URL.Path, "token") || strings.Contains(req.URL.Path, "auth") {
			log.Printf("Exiting token or auth")
			return resp, err
		}
		if resp.StatusCode == http.StatusTemporaryRedirect {
			log.Printf("Exiting StatusCode == StatusTemporaryRedirect")
			return resp, err
		}
		if resp.StatusCode != http.StatusUnauthorized {
			log.Printf("Exiting StatusCode != StatusUnauthorized or StatusTemporaryRedirect")
			log.Print(resp.StatusCode)
			return resp, err
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		accessToken, err := GetAuthResponse(t.account, t.username, t.password, t.scope)
		t.accessToken = accessToken
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %w", err)
		}
		//		log.Printf("Setting authorization header to %s\n", accessToken)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		resp.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		if accessToken != "" {
			_ = resp.Body.Close()
		}
	}
	return nil, fmt.Errorf("unauthorized - max attempts (%d)", 3)
}

func checkToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	if token.Valid {
		fmt.Println("You look nice today")
		return true
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return false
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			return false
		} else {
			return false
		}
	} else {
		return false
	}
}
