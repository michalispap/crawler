package main

import (
	"fmt"
)

func main() {
	normalizedURLExample, _ := normalizeURL("https://www.boot.dev/lessons/")
	fmt.Println(normalizedURLExample)
}
