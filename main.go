package main

import (
	"fmt"
	"os"
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

	fmt.Printf(getHTML(os.Args[1]))
}
