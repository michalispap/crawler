package main

import (
	"fmt"
	"net/url"
	"strings"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	normCurrent, _ := normalizeURL(rawCurrentURL)
	normBase, _ := normalizeURL(rawBaseURL)
	if !strings.Contains(normCurrent, normBase) {
		return
	}

	_, exists := pages[normCurrent]
	if exists {
		pages[normCurrent]++
		return
	}

	pages[normCurrent] = 1

	fmt.Printf("Currently scraping: %s\n", rawCurrentURL)

	rawHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	parsedURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	internalURLs, err := getURLsFromHTML(rawHTML, parsedURL)
	if err != nil {
		return
	}

	internalImages, err := getImagesFromHTML(rawHTML, parsedURL)
	if err != nil {
		return
	}

	totalLinks := append(internalURLs, internalImages...)

	for _, internalURL := range totalLinks {
		crawlPage(rawBaseURL, internalURL, pages)
	}
}
