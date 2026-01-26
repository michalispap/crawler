package main

import "testing"

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
