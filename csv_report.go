package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func writeCSVReport(pages map[string]PageData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	err = w.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return err
	}

	for _, v := range pages {
		err = w.Write([]string{v.URL, v.H1, v.FirstParagraph, strings.Join(v.OutgoingLinks, ";"), strings.Join(v.ImageURLs, ";")})
		if err != nil {
			return err
		}
	}

	return nil
}
