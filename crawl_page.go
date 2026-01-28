package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	normCurrent, _ := normalizeURL(rawCurrentURL)
	normBase, _ := normalizeURL(cfg.baseURL.String())
	if !strings.Contains(normCurrent, normBase) {
		return
	}

	if !cfg.addPageVisit(normCurrent) {
		return
	}

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
		cfg.wg.Add(1)
		go func(urlToCrawl string) {
			cfg.concurrencyControl <- struct{}{}
			defer func() {
				<-cfg.concurrencyControl
				cfg.wg.Done()
			}()
			cfg.crawlPage(urlToCrawl)
		}(internalURL)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, exists := cfg.pages[normalizedURL]; exists {
		return false
	}
	cfg.pages[normalizedURL] = PageData{}
	return true
}
