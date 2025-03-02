package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

// renderToString renders a templ.Component to a string for testing.
func renderToString(t *testing.T, component templ.Component) string {
	var buf bytes.Buffer

	err := component.Render(context.TODO(), &buf)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	return buf.String()
}

// parseHTML parses the HTML string into a goquery document for testing.
func parseHTML(t *testing.T, html string) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	return doc
}

func TestHomeContent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "HomeContent template structure and content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Render the HomeContent template
			html := renderToString(t, HomeContent())
			doc := parseHTML(t, html)

			// Test specific elements and attributes
			assert.Equal(t, "EUI-64 Calculator", doc.Find("h1.app-title").Text(), "Incorrect h1 text")
			assert.Equal(t, "Enter a MAC address and IPv6 prefix to calculate the EUI-64 address.", doc.Find("p.app-description").Text(), "Incorrect description text")
			assert.Equal(t, 1, doc.Find("form[hx-post='/calculate']").Length(), "Form not found")
			assert.Equal(t, 1, doc.Find("form[hx-target='.result-container']").Length(), "Form hx-target not found")
			assert.Equal(t, 1, doc.Find("form[hx-swap='innerHTML']").Length(), "Form hx-swap not found")
			assert.Equal(t, "MAC Address", doc.Find("label[for='mac']").Text(), "Incorrect MAC label text")
			assert.Equal(t, "xx-xx-xx-xx-xx-xx or xx:xx:xx:xx:xx:xx", doc.Find("input#mac").AttrOr("placeholder", ""), "Incorrect MAC input placeholder")
			assert.Equal(t, "[0-9a-fA-F]{2}([-:][0-9a-fA-F]{2}){5}", doc.Find("input#mac").AttrOr("pattern", ""), "Incorrect MAC pattern")
			assert.Equal(t, "MAC address must be in format xx-xx-xx-xx-xx-xx or xx:xx:xx:xx:xx:xx (e.g., 00-14-22-01-23-45 or 00:14:22:01:23:45)", doc.Find("input#mac").AttrOr("title", ""), "Incorrect MAC title")
			assert.Equal(t, "Start of IPv6 Address", doc.Find("label[for='ip-start']").Text(), "Incorrect IP label text")
			assert.Equal(t, "xxxx:xxxx:xxxx:xxxx", doc.Find("input#ip-start").AttrOr("placeholder", ""), "Incorrect IP input placeholder")
			assert.Equal(t, "Calculate", doc.Find("button.form-submit").Text(), "Incorrect submit button text")
			assert.Equal(t, "Clear", doc.Find("button.form-clear").Text(), "Incorrect reset button text")
			assert.Equal(t, 1, doc.Find("div.form-results.hidden").Length(), "Form results div not found or not hidden")
			assert.Equal(t, 1, doc.Find("div.result-container.hidden").Length(), "Result container div not found or not hidden")

			// Check that unwanted content is not present
			assert.Equal(t, 0, doc.Find(".unexpected-class").Length(), "Unexpected class found in HomeContent template")
		})
	}
}

func TestHome(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Home template structure and content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Render the Home template (which includes Layout and HomeContent)
			html := renderToString(t, Home())
			doc := parseHTML(t, html)

			// Test specific elements and attributes in the full HTML structure
			assert.Equal(t, "en", doc.Find("html").AttrOr("lang", ""), "Incorrect html lang attribute") // Fixed expected value
			assert.Equal(t, "EUI-64 Calculator", doc.Find("title").Text(), "Incorrect title")
			assert.Equal(t, 1, doc.Find("meta[charset='UTF-8']").Length(), "Meta charset not found")
			assert.Equal(t, 1, doc.Find("meta[name='viewport']").Length(), "Meta viewport not found")
			assert.Equal(t, "width=device-width, initial-scale=1.0", doc.Find("meta[name='viewport']").AttrOr("content", ""), "Incorrect viewport content")
			assert.Equal(t, 1, doc.Find("link[href='/static/favicon.ico']").Length(), "Favicon link not found")
			assert.Equal(t, "image/x-icon", doc.Find("link[href='/static/favicon.ico']").AttrOr("type", ""), "Incorrect favicon type")
			assert.Equal(t, 1, doc.Find("link[href='/static/styles.css']").Length(), "Stylesheet link not found")
			assert.Equal(t, 1, doc.Find("link[href='https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap']").Length(), "Google Fonts link not found")
			assert.Equal(t, 1, doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").Length(), "HTMX script not found")
			assert.Equal(t, "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+", doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").AttrOr("integrity", ""), "Incorrect HTMX integrity")
			assert.Equal(t, "anonymous", doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").AttrOr("crossorigin", ""), "Incorrect HTMX crossorigin")
			assert.Equal(t, 1, doc.Find("div.app-container").Length(), "App container div not found")
			assert.Equal(t, "EUI-64 Calculator", doc.Find("h1.app-title").Text(), "Incorrect h1 text")
			assert.Equal(t, 1, doc.Find("script").FilterFunction(func(i int, s *goquery.Selection) bool {
				return strings.Contains(s.Text(), "document.body.addEventListener('htmx:afterSwap'")
			}).Length(), "HTMX afterSwap script not found")

			// Check that unwanted content is not present
			assert.Equal(t, 0, doc.Find(".unexpected-class").Length(), "Unexpected class found in Home template")
		})
	}
}

func TestLayout(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		content templ.Component
	}{
		{
			name:    "Layout template structure and content",
			title:   "Test Title",
			content: HomeContent(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Render the Layout template with the given title and content
			html := renderToString(t, Layout(tt.title, tt.content))
			doc := parseHTML(t, html)

			// Test specific elements and attributes in the full HTML structure
			assert.Equal(t, "en", doc.Find("html").AttrOr("lang", ""), "Incorrect html lang attribute") // Fixed expected value
			assert.Equal(t, tt.title, doc.Find("title").Text(), "Incorrect title")
			assert.Equal(t, 1, doc.Find("meta[charset='UTF-8']").Length(), "Meta charset not found")
			assert.Equal(t, 1, doc.Find("meta[name='viewport']").Length(), "Meta viewport not found")
			assert.Equal(t, "width=device-width, initial-scale=1.0", doc.Find("meta[name='viewport']").AttrOr("content", ""), "Incorrect viewport content")
			assert.Equal(t, 1, doc.Find("link[href='/static/favicon.ico']").Length(), "Favicon link not found")
			assert.Equal(t, "image/x-icon", doc.Find("link[href='/static/favicon.ico']").AttrOr("type", ""), "Incorrect favicon type")
			assert.Equal(t, 1, doc.Find("link[href='/static/styles.css']").Length(), "Stylesheet link not found")
			assert.Equal(t, 1, doc.Find("link[href='https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap']").Length(), "Google Fonts link not found")
			assert.Equal(t, 1, doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").Length(), "HTMX script not found")
			assert.Equal(t, "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+", doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").AttrOr("integrity", ""), "Incorrect HTMX integrity")
			assert.Equal(t, "anonymous", doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").AttrOr("crossorigin", ""), "Incorrect HTMX crossorigin")
			assert.Equal(t, 1, doc.Find("div.app-container").Length(), "App container div not found")
			assert.Equal(t, "EUI-64 Calculator", doc.Find("h1.app-title").Text(), "Incorrect h1 text")
			assert.Equal(t, 1, doc.Find("script").FilterFunction(func(i int, s *goquery.Selection) bool {
				return strings.Contains(s.Text(), "document.body.addEventListener('htmx:afterSwap'")
			}).Length(), "HTMX afterSwap script not found")

			// Check that unwanted content is not present
			assert.Equal(t, 0, doc.Find(".unexpected-class").Length(), "Unexpected class found in Layout template")
		})
	}
}

func TestResult(t *testing.T) {
	tests := []struct {
		name      string
		data      ResultData
		assertDoc func(t *testing.T, doc *goquery.Document)
	}{
		{
			name: "Result template with success data",
			data: ResultData{
				InterfaceID: "0214:22ff:fe01:2345",
				FullIP:      "2001:0db8:85a3:0000:0214:22ff:fe01:2345",
				Error:       "",
			},
			assertDoc: func(t *testing.T, doc *goquery.Document) {
				assert.Equal(t, "End of IPv6 Address", doc.Find("label[for='ip']").Text(), "Incorrect interface ID label text")
				assert.Equal(t, "0214:22ff:fe01:2345", doc.Find("input[readonly]").First().AttrOr("value", ""), "Incorrect interface ID value")
				assert.Equal(t, "IPv6 Address", doc.Find("label[for='ip-full']").Text(), "Incorrect full IP label text")
				assert.Equal(t, "2001:0db8:85a3:0000:0214:22ff:fe01:2345", doc.Find("input[readonly]").Last().AttrOr("value", ""), "Incorrect full IP value")
				assert.Equal(t, 0, doc.Find("p.error-message").Length(), "Error message should not be present")
			},
		},
		{
			name: "Result template with error data",
			data: ResultData{
				InterfaceID: "",
				FullIP:      "",
				Error:       "Invalid MAC address",
			},
			assertDoc: func(t *testing.T, doc *goquery.Document) {
				assert.Equal(t, "Invalid MAC address", doc.Find("p.error-message").Text(), "Incorrect error message")
				assert.Equal(t, 0, doc.Find("label[for='ip']").Length(), "Success fields should not be present")
				assert.Equal(t, 0, doc.Find("input[readonly]").Length(), "Success fields should not be present")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Render the Result template with the given data
			html := renderToString(t, Result(tt.data))
			doc := parseHTML(t, html)

			// Run the assertions specific to this test case
			tt.assertDoc(t, doc)
		})
	}
}
