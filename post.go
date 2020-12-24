package main

import (
	"encoding/base64"
	"net/http"
)

type gistFile struct {
	Type    string `json:"type,omitempty"`
	Content string `json:"content"`
}

type postRequest struct {
	Description string              `json:"description,omitempty"`
	Files       map[string]gistFile `json:"files"`
	Public      bool                `json:"public"`
}

type postResponse struct {
	HTMLURL string `json:"html_url"`
	URL     string `json:"url,omitempty"`
	ID      string `json:"id"`
}

func post(description string, files map[string][]byte) (*postResponse, error) {
	header := http.Header{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/json"},
	}
	if args.AccessToken != "" {
		header.Set("Authorization", "Bearer "+args.AccessToken)
	}
	postReq := postRequest{description, make(map[string]gistFile), !args.Private}
	for k, v := range files {
		content := gistFile{Content: string(v)}
		if args.Base64 {
			content.Content = base64.StdEncoding.EncodeToString(v)
			content.Type = "application/base64"
		}
		postReq.Files[k] = content
	}
	var postResp postResponse
	err := rest("POST", buildURL("api", "/gists"), header, &postReq, &postResp)
	if err != nil {
		return nil, err
	}
	return &postResp, nil
}
