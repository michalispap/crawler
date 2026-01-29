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
	maxPages           int
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}

	normCurrent, _ := normalizeURL(rawCurrentURL)
	normBase, _ := normalizeURL(cfg.baseURL.String())
	if !strings.Contains(normCurrent, normBase) {
		cfg.mu.Unlock()
		return
	}

	if !cfg.addPageVisit(normCurrent) {
		cfg.mu.Unlock()
		return
	}

	cfg.mu.Unlock()

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
	if _, exists := cfg.pages[normalizedURL]; exists {
		return false
	}
	cfg.pages[normalizedURL] = PageData{}
	return true
}
