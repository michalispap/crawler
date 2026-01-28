package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("no website provided")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		fmt.Printf("too many arguments provided")
		os.Exit(1)
	} else {
		fmt.Printf("starting crawl of: %s\n", os.Args[1])
	}

	website := os.Args[1]

	parsedBaseURL, err := url.Parse(website)
	if err != nil {
		fmt.Printf("invalid URL: %v\n", err)
		os.Exit(1)
	}

	cfg := config{
		pages:              make(map[string]PageData),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 1),
		wg:                 &sync.WaitGroup{},
	}

	cfg.crawlPage(website)

	cfg.wg.Wait()

	for page, data := range cfg.pages {
		fmt.Printf("%s: %+v\n", page, data)
	}
}
