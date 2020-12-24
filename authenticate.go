package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

type loginRequest struct {
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
}

type loginResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type accessTokenRequest struct {
	ClientID   string `json:"client_id"`
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func loadToken() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	tokenPath := filepath.Join(usr.HomeDir, ".config", "go-gist", "token.txt")
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return err
	}
	args.AccessToken = string(token)
	return nil
}

func saveToken() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}
	path := filepath.Join(usr.HomeDir, ".config", "go-gist")
	err = os.MkdirAll(path, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create config directory\n")
	}
	err = ioutil.WriteFile(filepath.Join(path, "token.txt"), []byte(args.AccessToken), 0600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot write token file to disk\n")
	}
	return nil
}

func authenticate() error {
	if args.AccessToken != "" {
		return nil
	}
	header := http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	loginReq := loginRequest{args.ClientID, "gist"}
	var loginResp loginResponse
	err := rest("POST", buildURL("", "/login/device/code"), header, &loginReq, &loginResp)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Go to the following URL:\n%s\nEnter: %s to verify login\n", loginResp.VerificationURI, loginResp.UserCode)
	accessTokenReq := accessTokenRequest{
		args.ClientID, loginResp.DeviceCode,
		"urn:ietf:params:oauth:grant-type:device_code",
	}
	var accessTokenResp accessTokenResponse
	secondsElapsed := 0
	for true {
		if secondsElapsed >= loginResp.ExpiresIn || (args.Timeout > 0 && secondsElapsed >= args.Timeout) {
			return errors.New("timed out")
		}
		time.Sleep(time.Duration(loginResp.Interval) * time.Second)
		secondsElapsed += loginResp.Interval
		err = rest("POST", buildURL("", "/login/oauth/access_token"), header, &accessTokenReq, &accessTokenResp)
		if githubErr, ok := err.(*githubError); (ok && githubErr.StatusCode/100 != 2) || (!ok && err != nil) {
			return err
		}
		if accessTokenResp.AccessToken != "" {
			args.AccessToken = accessTokenResp.AccessToken
			saveToken()
			return nil
		}
	}
	return nil
}
