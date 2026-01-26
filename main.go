package main

import (
	"fmt"
)

func main() {
	normalizedURLExample, _ := normalizeURL("https://www.boot.dev/courses/")
	fmt.Println(getH1FromHTML(normalizedURLExample))
}
