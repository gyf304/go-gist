package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type githubError struct {
	ErrorStr         string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorURI         string `json:"error_uri"`
	Message          string `json:"message"`
	StatusCode       int
}

func (e *githubError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.ErrorDescription
}

func mustParseURL(rawurl string) *url.URL {
	result, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	return result
}

func buildURL(subdomain string, path ...string) *url.URL {
	gurl, err := url.Parse(args.GithubURL)
	if err != nil {
		panic(err)
	}
	if subdomain != "" {
		gurl.Host = subdomain + "." + gurl.Host
	}
	gurl.Path = strings.Join(path, "/")
	return gurl
}

func rest(method string, url *url.URL, header http.Header, req interface{}, resp interface{}) error {
	var reqBuf []byte
	var err error
	if req != nil {
		reqBuf, err = json.Marshal(req)
		if err != nil {
			return err
		}
	}
	httpResp, err := http.DefaultClient.Do(&http.Request{
		Method: method,
		URL:    url,
		Header: header,
		Body:   ioutil.NopCloser(bytes.NewReader(reqBuf)),
	})
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	var githubErr githubError
	_ = json.Unmarshal(body, &githubErr)
	if githubErr.ErrorStr != "" || httpResp.StatusCode/100 != 2 {
		githubErr.StatusCode = httpResp.StatusCode
		return &githubErr
	}
	return json.Unmarshal(body, resp)
}
