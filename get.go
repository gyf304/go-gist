package main

import (
	"encoding/base64"
	"net/http"
)

type getResponse struct {
	URL   string              `json:"url,omitempty"`
	Files map[string]gistFile `json:"files"`
}

func get(id string) (files map[string][]byte, err error) {
	header := http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + args.AccessToken},
	}
	var getResp getResponse
	err = rest("GET", buildURL("api", "/gists", id), header, nil, &getResp)
	if err != nil {
		return nil, err
	}
	files = make(map[string][]byte)
	for k, v := range getResp.Files {
		content := []byte(v.Content)
		if args.Base64 || v.Type == "application/base64" {
			content, err = base64.StdEncoding.DecodeString(v.Content)
			if err != nil {
				return nil, err
			}
		}
		files[k] = content
	}
	return files, nil
}
