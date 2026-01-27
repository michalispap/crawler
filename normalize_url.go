package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	urlStruct, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	urlStruct.Scheme = ""
	urlStruct.Path, _ = strings.CutSuffix(urlStruct.Path, "/")

	return urlStruct.Host + urlStruct.Path, nil
}
