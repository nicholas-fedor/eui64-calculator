package ui_test

import (
	"bytes"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/a-h/templ"
	"github.com/nicholas-fedor/eui64-calculator/ui"
	"github.com/stretchr/testify/require"
)

const (
	viewportContent = "width=device-width, initial-scale=1.0"
	googleFontsURL  = "https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap"
	htmxURL         = "https://unpkg.com/htmx.org@2.0.4"
)

func renderToString(t *testing.T, component templ.Component) string {
	t.Helper()

	var buffer bytes.Buffer

	err := component.Render(t.Context(), &buffer)
	require.NoError(t, err, "Failed to render component")

	return buffer.String()
}

func parseHTML(t *testing.T, html string) *goquery.Document {
	t.Helper()

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
	require.NoError(t, err, "Failed to parse HTML")

	return doc
}

func TestHomeContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "Home content check"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := renderToString(t, ui.Home())
			doc := parseHTML(t, html)

			expectedDesc := "Enter a MAC address and IPv6 prefix to calculate the EUI-64 address."
			actualDesc := doc.Find("p.app-description").Text()
			require.Equal(t, expectedDesc, actualDesc, "Incorrect description text")

			require.Equal(t, "xx-xx-xx-xx-xx-xx", doc.Find("input#mac").
				AttrOr("placeholder", ""), "Incorrect MAC input placeholder")

			require.Equal(t, "xxxx:xxxx:xxxx:xxxx", doc.Find("input#ip-start").
				AttrOr("placeholder", ""), "Incorrect IP input placeholder")

			require.Equal(t, 1, doc.Find("button[type='submit']").Length(), "Submit button not found")
		})
	}
}

func TestHome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{name: "Home page layout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := renderToString(t, ui.Home())
			doc := parseHTML(t, html)

			require.Equal(t, "EUI-64 Calculator", doc.Find("title").Text(), "Incorrect title")

			require.Equal(t, viewportContent, doc.Find("meta[name='viewport']").
				AttrOr("content", ""), "Incorrect viewport content")

			require.Equal(t, "image/x-icon", doc.Find("link[href='/static/favicon.ico']").
				AttrOr("type", ""), "Incorrect favicon type")

			require.Equal(t, 1, doc.Find("link[href='"+googleFontsURL+"']").
				Length(), "Google Fonts link not found")

			require.Equal(t, "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+",
				doc.Find("script[src='"+htmxURL+"']").
					AttrOr("integrity", ""), "Incorrect HTMX integrity")

			require.Equal(t, "anonymous", doc.Find("script[src='"+htmxURL+"']").
				AttrOr("crossorigin", ""), "Incorrect HTMX crossorigin")

			require.Equal(t, 1, doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
				return s.AttrOr("src", "") == htmxURL
			}).Length(), "HTMX script not found")
		})
	}
}

// TestResultValid tests the Result function with valid data.
func TestResultValid(t *testing.T) {
	t.Parallel()

	data := ui.ResultData{
		InterfaceID: "0214:22ff:fe01:2345",
		FullIP:      "2001:0db8:85a3:0000:0214:22ff:fe01:2345",
		Error:       "",
	}

	html := renderToString(t, ui.Result(data))
	doc := parseHTML(t, html)

	assertDocValid(t, doc, data)
}

// TestResultError tests the Result function with an error.
func TestResultError(t *testing.T) {
	t.Parallel()

	data := ui.ResultData{
		InterfaceID: "",
		FullIP:      "",
		Error:       "Invalid input",
	}

	html := renderToString(t, ui.Result(data))
	doc := parseHTML(t, html)

	assertDocError(t, doc, data.Error) // Pass only the Error field
}

func assertDocValid(t *testing.T, doc *goquery.Document, _ ui.ResultData) {
	t.Helper()

	require.Equal(t, "EUI-64 Calculator", doc.Find("title").Text(), "Incorrect title")

	require.Equal(t, viewportContent, doc.Find("meta[name='viewport']").
		AttrOr("content", ""), "Incorrect viewport content")

	require.Equal(t, "image/x-icon", doc.Find("link[href='/static/favicon.ico']").
		AttrOr("type", ""), "Incorrect favicon type")

	require.Equal(t, 1, doc.Find("link[href='"+googleFontsURL+"']").
		Length(), "Google Fonts link not found")

	require.Equal(t, "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+",
		doc.Find("script[src='"+htmxURL+"']").
			AttrOr("integrity", ""), "Incorrect HTMX integrity")

	require.Equal(t, "anonymous", doc.Find("script[src='"+htmxURL+"']").
		AttrOr("crossorigin", ""), "Incorrect HTMX crossorigin")

	require.Equal(t, 1, doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.AttrOr("src", "") == htmxURL
	}).Length(), "HTMX script not found")

	require.Equal(t, "0214:22ff:fe01:2345", doc.Find("input[readonly]").
		First().AttrOr("value", ""), "Incorrect interface ID value")
	require.Equal(t, "2001:0db8:85a3:0000:0214:22ff:fe01:2345",
		doc.Find("input[readonly]").Last().AttrOr("value", ""), "Incorrect full IP value")
}

func assertDocError(t *testing.T, doc *goquery.Document, errorMsg string) {
	t.Helper()

	require.Equal(t, "EUI-64 Calculator", doc.Find("title").Text(), "Incorrect title")

	require.Equal(t, viewportContent, doc.Find("meta[name='viewport']").
		AttrOr("content", ""), "Incorrect viewport content")

	require.Equal(t, "image/x-icon", doc.Find("link[href='/static/favicon.ico']").
		AttrOr("type", ""), "Incorrect favicon type")

	require.Equal(t, 1, doc.Find("link[href='"+googleFontsURL+"']").
		Length(), "Google Fonts link not found")

	require.Equal(t, "sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+",
		doc.Find("script[src='"+htmxURL+"']").
			AttrOr("integrity", ""), "Incorrect HTMX integrity")

	require.Equal(t, "anonymous", doc.Find("script[src='"+htmxURL+"']").
		AttrOr("crossorigin", ""), "Incorrect HTMX crossorigin")

	require.Equal(t, 1, doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return s.AttrOr("src", "") == htmxURL
	}).Length(), "HTMX script not found")

	require.Contains(t, doc.Find("p.error-message").Text(), errorMsg, "Error message not found")
}
