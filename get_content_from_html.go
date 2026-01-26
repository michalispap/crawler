package main

import (
	"net/url"
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

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	allURLs := make([]string, 0)
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		parsedURL, _ := baseURL.Parse(href)
		allURLs = append(allURLs, parsedURL.String())
	})

	return allURLs, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	allURLs := make([]string, 0)
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		parsedURL, _ := baseURL.Parse(src)
		allURLs = append(allURLs, parsedURL.String())
	})

	return allURLs, nil
}

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) PageData {
	var outgoingLinks []string
	var imageURLs []string
	parsedURL, err := url.Parse(pageURL)
	if err == nil {
		outgoingLinks, _ = getURLsFromHTML(html, parsedURL)
		imageURLs, _ = getImagesFromHTML(html, parsedURL)
	}
	return PageData{
		URL:            pageURL,
		H1:             getH1FromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
