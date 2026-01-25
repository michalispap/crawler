package main

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	urlStruct, err := url.Parse(rawURL) // Parse raw URL
	if err != nil {
		return "", err
	}

	// Normalize
	urlStruct.Scheme = ""
	urlStruct.Path, _ = strings.CutSuffix(urlStruct.Path, "/")

	return urlStruct.Host + urlStruct.Path, nil
}
