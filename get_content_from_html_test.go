package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic h1",
			input:    "<html><body><h1>Test Title</h1></body></html>",
			expected: "Test Title",
		},
		{
			name:     "multiple h1 use first",
			input:    "<h1>First</h1><h1>Second</h1>",
			expected: "First",
		},
		{
			name:     "no h1",
			input:    "<html><body><p>No heading</p></body></html>",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getH1FromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected %q, got %q", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "main paragraph has priority",
			input: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<p>Main paragraph.</p>
				</main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "no main use first body paragraph",
			input: `<html><body>
				<p>First paragraph.</p>
				<p>Second paragraph.</p>
			</body></html>`,
			expected: "First paragraph.",
		},
		{
			name:     "no paragraph",
			input:    "<html><body><div>No content</div></body></html>",
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected %q, got %q", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputHTML string
		expected  []string
	}{
		{
			name:      "single absolute url",
			inputURL:  "https://blog.boot.dev",
			inputHTML: `<html><body><a href="https://blog.boot.dev">Boot</a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		},
		{
			name:      "single relative url",
			inputURL:  "https://blog.boot.dev",
			inputHTML: `<html><body><a href="/path">Path</a></body></html>`,
			expected:  []string{"https://blog.boot.dev/path"},
		},
		{
			name:     "multiple links mixed",
			inputURL: "https://blog.boot.dev",
			inputHTML: `<html><body>
				<a href="/one">one</a>
				<a href="https://example.com/two">two</a>
			</body></html>`,
			expected: []string{
				"https://blog.boot.dev/one",
				"https://example.com/two",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Fatalf("Test %v - %s FAIL: could not parse base url: %v", i, tc.name, err)
			}

			actual, err := getURLsFromHTML(tc.inputHTML, baseURL)
			if err != nil {
				t.Fatalf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf(
					"Test %v - %s FAIL:\nexpected %v\ngot      %v",
					i,
					tc.name,
					tc.expected,
					actual,
				)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "single absolute image src",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><img src="https://cdn.example.com/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://cdn.example.com/logo.png"},
		},
		{
			name:      "single relative image src",
			inputURL:  "https://blog.boot.dev",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://blog.boot.dev/logo.png"},
		},
		{
			name:     "missing src is ignored",
			inputURL: "https://blog.boot.dev",
			inputBody: `<html><body>
				<img alt="No src here">
				<img src="/ok.png" alt="Ok">
			</body></html>`,
			expected: []string{"https://blog.boot.dev/ok.png"},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Fatalf("Test %v - %s FAIL: couldn't parse input URL: %v", i, tc.name, err)
			}

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected %v, got %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:             "https://blog.boot.dev",
		H1:              "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
		ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}