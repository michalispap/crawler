package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	h1 := doc.Find("h1").First()
	if h1.Length() == 0 {
		return ""
	}

	return h1.Text()
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	if pInMain := doc.Find("main").Find("p"); pInMain.Length() != 0 {
		return pInMain.First().Text()
	}

	p := doc.Find("p").First()
	if p.Length() == 0 {
		return ""
	}

	return p.Text()
}
