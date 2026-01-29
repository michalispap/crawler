package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

// Usage: ./crawler URL maxConcurrency maxPages

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("no website provided")
		os.Exit(1)
	} else if len(os.Args) > 4 {
		fmt.Printf("too many arguments provided")
		os.Exit(1)
	} else {
		fmt.Printf("starting crawl of: %s\nmax concurrency: %s\nmax pages: %s\n", os.Args[1], os.Args[2], os.Args[3])
	}

	website := os.Args[1]

	parsedBaseURL, err := url.Parse(website)
	if err != nil {
		fmt.Printf("invalid URL: %v\n", err)
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("invalid integer:", err)
		return
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("invalid integer:", err)
		return
	}

	cfg := config{
		pages:              make(map[string]PageData),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.crawlPage(website)

	cfg.wg.Wait()

	for page, data := range cfg.pages {
		fmt.Printf("%s: %+v\n", page, data)
	}
}
