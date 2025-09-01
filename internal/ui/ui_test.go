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
	t.Helper()

	var buf bytes.Buffer

	err := component.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	return buf.String()
}

// parseHTML parses the HTML string into a goquery document for testing.
func parseHTML(t *testing.T, html string) *goquery.Document {
	t.Helper()

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
			html := renderToString(t, HomeContent())
			doc := parseHTML(t, html)

			assert.Equal(
				t,
				"EUI-64 Calculator",
				doc.Find("h1.app-title").Text(),
				"Incorrect h1 text",
			)
			assert.Equal(
				t,
				"Enter a MAC address and IPv6 prefix to calculate the EUI-64 address.",
				doc.Find("p.app-description").Text(),
				"Incorrect description text",
			)
			assert.Equal(t, 1, doc.Find("form[hx-post='/calculate']").Length(), "Form not found")
			assert.Equal(
				t,
				1,
				doc.Find("form[hx-target='.result-container']").Length(),
				"Form hx-target not found",
			)
			assert.Equal(
				t,
				1,
				doc.Find("form[hx-swap='innerHTML']").Length(),
				"Form hx-swap not found",
			)
			assert.Equal(
				t,
				"MAC Address",
				doc.Find("label[for='mac']").Text(),
				"Incorrect MAC label text",
			)
			assert.Equal(
				t,
				"xx-xx-xx-xx-xx-xx or xx:xx:xx:xx:xx:xx",
				doc.Find("input#mac").AttrOr("placeholder", ""),
				"Incorrect MAC input placeholder",
			)
			assert.Equal(
				t,
				"[0-9a-fA-F]{2}([-:][0-9a-fA-F]{2}){5}",
				doc.Find("input#mac").AttrOr("pattern", ""),
				"Incorrect MAC pattern",
			)
			assert.Equal(
				t,
				"MAC address must be in format xx-xx-xx-xx-xx-xx or xx:xx:xx:xx:xx:xx (e.g., 00-14-22-01-23-45 or 00:14:22:01:23:45)",
				doc.Find("input#mac").AttrOr("title", ""),
				"Incorrect MAC title",
			)
			assert.Equal(
				t,
				"mac-copy",
				doc.Find("input#mac").AttrOr("aria-describedby", ""),
				"Incorrect MAC aria-describedby",
			)
			assert.Equal(
				t,
				"Start of IPv6 Address",
				doc.Find("label[for='ip-start']").Text(),
				"Incorrect IP label text",
			)
			assert.Equal(
				t,
				"xxxx:xxxx:xxxx:xxxx",
				doc.Find("input#ip-start").AttrOr("placeholder", ""),
				"Incorrect IP input placeholder",
			)
			assert.Equal(
				t,
				"ip-start-copy",
				doc.Find("input#ip-start").AttrOr("aria-describedby", ""),
				"Incorrect IP aria-describedby",
			)
			assert.Equal(
				t,
				"Calculate",
				doc.Find("button.form-submit").Text(),
				"Incorrect submit button text",
			)
			assert.Equal(
				t,
				"Clear",
				doc.Find("button.form-clear").Text(),
				"Incorrect reset button text",
			)
			assert.Equal(
				t,
				1,
				doc.Find("div.form-results.hidden").Length(),
				"Form results div not found or not hidden",
			)
			assert.Equal(
				t,
				1,
				doc.Find("div.result-container.hidden").Length(),
				"Result container div not found or not hidden",
			)

			// Test copy buttons for input fields
			macCopyBtn := doc.Find("#copy-mac")
			assert.Equal(t, 1, macCopyBtn.Length(), "MAC copy button not found")
			assert.Equal(
				t,
				"button",
				macCopyBtn.AttrOr("type", ""),
				"MAC copy button should have type='button'",
			)
			assert.Equal(
				t,
				"Copy MAC Address",
				macCopyBtn.AttrOr("aria-label", ""),
				"Incorrect MAC copy button aria-label",
			)
			assert.Equal(t, 1, macCopyBtn.Find(".copy-icon").Length(), "MAC copy icon not found")
			assert.Equal(
				t,
				1,
				macCopyBtn.Find(".copy-tooltip").Length(),
				"MAC copy tooltip not found",
			)

			ipStartCopyBtn := doc.Find("#copy-ip-start")
			assert.Equal(t, 1, ipStartCopyBtn.Length(), "IPv6 Prefix copy button not found")
			assert.Equal(
				t,
				"button",
				ipStartCopyBtn.AttrOr("type", ""),
				"IPv6 Prefix copy button should have type='button'",
			)
			assert.Equal(
				t,
				"Copy IPv6 Prefix",
				ipStartCopyBtn.AttrOr("aria-label", ""),
				"Incorrect IPv6 Prefix copy button aria-label",
			)
			assert.Equal(
				t,
				1,
				ipStartCopyBtn.Find(".copy-icon").Length(),
				"IPv6 Prefix copy icon not found",
			)
			assert.Equal(
				t,
				1,
				ipStartCopyBtn.Find(".copy-tooltip").Length(),
				"IPv6 Prefix copy tooltip not found",
			)

			// Verify copy button event listeners
			assert.Equal(
				t,
				1,
				doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
					return strings.Contains(
						s.Text(),
						`document.getElementById("copy-mac").addEventListener("click"`,
					)
				}).Length(),
				"MAC copy button event listener not found",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
					return strings.Contains(
						s.Text(),
						`document.getElementById("copy-ip-start").addEventListener("click"`,
					)
				}).Length(),
				"IPv6 Prefix copy button event listener not found",
			)

			// Check that unwanted content is not present
			assert.Equal(
				t,
				0,
				doc.Find(".unexpected-class").Length(),
				"Unexpected class found in HomeContent template",
			)
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
			html := renderToString(t, Home())
			doc := parseHTML(t, html)

			assert.Equal(
				t,
				"en",
				doc.Find("html").AttrOr("lang", ""),
				"Incorrect html lang attribute",
			)
			assert.Equal(t, "EUI-64 Calculator", doc.Find("title").Text(), "Incorrect title")
			assert.Equal(t, 1, doc.Find("meta[charset='UTF-8']").Length(), "Meta charset not found")
			assert.Equal(
				t,
				1,
				doc.Find("meta[name='viewport']").Length(),
				"Meta viewport not found",
			)
			assert.Equal(
				t,
				"width=device-width, initial-scale=1.0",
				doc.Find("meta[name='viewport']").AttrOr("content", ""),
				"Incorrect viewport content",
			)
			assert.Equal(
				t,
				1,
				doc.Find("link[href='/static/favicon.ico']").Length(),
				"Favicon link not found",
			)
			assert.Equal(
				t,
				"image/x-icon",
				doc.Find("link[href='/static/favicon.ico']").AttrOr("type", ""),
				"Incorrect favicon type",
			)
			assert.Equal(
				t,
				1,
				doc.Find("link[href='/static/styles.css']").Length(),
				"Stylesheet link not found",
			)
			assert.Equal(
				t,
				0,
				doc.Find("link[href*='fonts.googleapis.com']").Length(),
				"Google Fonts link should not be present",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").Length(),
				"HTMX script not found",
			)
			assert.Equal(
				t,
				"sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+",
				doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").AttrOr("integrity", ""),
				"Incorrect HTMX integrity",
			)
			assert.Equal(
				t,
				"anonymous",
				doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").
					AttrOr("crossorigin", ""),
				"Incorrect HTMX crossorigin",
			)
			assert.Equal(
				t,
				1,
				doc.Find("div.app-container").Length(),
				"App container div not found",
			)
			assert.Equal(
				t,
				"EUI-64 Calculator",
				doc.Find("h1.app-title").Text(),
				"Incorrect h1 text",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
					return strings.Contains(
						s.Text(),
						"document.body.addEventListener('htmx:afterSwap'",
					)
				}).Length(),
				"HTMX afterSwap script not found",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
					return strings.Contains(
						s.Text(),
						"function copyToClipboard(elementId, buttonId)",
					)
				}).Length(),
				"copyToClipboard script not found",
			)
			assert.Equal(
				t,
				0,
				doc.Find(".unexpected-class").Length(),
				"Unexpected class found in Home template",
			)
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
			html := renderToString(t, Layout(tt.title, tt.content))
			doc := parseHTML(t, html)

			assert.Equal(
				t,
				"en",
				doc.Find("html").AttrOr("lang", ""),
				"Incorrect html lang attribute",
			)
			assert.Equal(t, tt.title, doc.Find("title").Text(), "Incorrect title")
			assert.Equal(t, 1, doc.Find("meta[charset='UTF-8']").Length(), "Meta charset not found")
			assert.Equal(
				t,
				1,
				doc.Find("meta[name='viewport']").Length(),
				"Meta viewport not found",
			)
			assert.Equal(
				t,
				"width=device-width, initial-scale=1.0",
				doc.Find("meta[name='viewport']").AttrOr("content", ""),
				"Incorrect viewport content",
			)
			assert.Equal(
				t,
				1,
				doc.Find("link[href='/static/favicon.ico']").Length(),
				"Favicon link not found",
			)
			assert.Equal(
				t,
				"image/x-icon",
				doc.Find("link[href='/static/favicon.ico']").AttrOr("type", ""),
				"Incorrect favicon type",
			)
			assert.Equal(
				t,
				1,
				doc.Find("link[href='/static/styles.css']").Length(),
				"Stylesheet link not found",
			)
			assert.Equal(
				t,
				0,
				doc.Find("link[href*='fonts.googleapis.com']").Length(),
				"Google Fonts link should not be present",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").Length(),
				"HTMX script not found",
			)
			assert.Equal(
				t,
				"sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+",
				doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").AttrOr("integrity", ""),
				"Incorrect HTMX integrity",
			)
			assert.Equal(
				t,
				"anonymous",
				doc.Find("script[src='https://unpkg.com/htmx.org@2.0.4']").
					AttrOr("crossorigin", ""),
				"Incorrect HTMX crossorigin",
			)
			assert.Equal(
				t,
				1,
				doc.Find("div.app-container").Length(),
				"App container div not found",
			)
			assert.Equal(
				t,
				"EUI-64 Calculator",
				doc.Find("h1.app-title").Text(),
				"Incorrect h1 text",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
					return strings.Contains(
						s.Text(),
						"document.body.addEventListener('htmx:afterSwap'",
					)
				}).Length(),
				"HTMX afterSwap script not found",
			)
			assert.Equal(
				t,
				1,
				doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
					return strings.Contains(
						s.Text(),
						"function copyToClipboard(elementId, buttonId)",
					)
				}).Length(),
				"copyToClipboard script not found",
			)
			assert.Equal(
				t,
				0,
				doc.Find(".unexpected-class").Length(),
				"Unexpected class found in Layout template",
			)
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
				t.Helper()
				assert.Equal(
					t,
					"End of IPv6 Address",
					doc.Find("label[for='interface-id']").Text(),
					"Incorrect interface ID label text",
				)
				assert.Equal(
					t,
					"0214:22ff:fe01:2345",
					doc.Find("input#interface-id").AttrOr("value", ""),
					"Incorrect interface ID value",
				)
				assert.Equal(
					t,
					"interface-id-copy",
					doc.Find("input#interface-id").AttrOr("aria-describedby", ""),
					"Incorrect interface ID aria-describedby",
				)
				assert.Equal(
					t,
					"IPv6 Address",
					doc.Find("label[for='ip-full']").Text(),
					"Incorrect full IP label text",
				)
				assert.Equal(
					t,
					"2001:0db8:85a3:0000:0214:22ff:fe01:2345",
					doc.Find("input#ip-full").AttrOr("value", ""),
					"Incorrect full IP value",
				)
				assert.Equal(
					t,
					"ip-full-copy",
					doc.Find("input#ip-full").AttrOr("aria-describedby", ""),
					"Incorrect full IP aria-describedby",
				)
				assert.Equal(
					t,
					0,
					doc.Find("p.error-message").Length(),
					"Error message should not be present",
				)

				// Test copy buttons for result fields
				interfaceCopyBtn := doc.Find("#copy-interface")
				assert.Equal(t, 1, interfaceCopyBtn.Length(), "Interface ID copy button not found")
				assert.Equal(
					t,
					"Copy Interface ID",
					interfaceCopyBtn.AttrOr("aria-label", ""),
					"Incorrect Interface ID copy button aria-label",
				)
				assert.Equal(
					t,
					1,
					interfaceCopyBtn.Find(".copy-icon").Length(),
					"Interface ID copy icon not found",
				)
				assert.Equal(
					t,
					1,
					interfaceCopyBtn.Find(".copy-tooltip").Length(),
					"Interface ID copy tooltip not found",
				)

				fullIPCopyBtn := doc.Find("#copy-ip-full")
				assert.Equal(t, 1, fullIPCopyBtn.Length(), "IPv6 Address copy button not found")
				assert.Equal(
					t,
					"Copy IPv6 Address",
					fullIPCopyBtn.AttrOr("aria-label", ""),
					"Incorrect IPv6 Address copy button aria-label",
				)
				assert.Equal(
					t,
					1,
					fullIPCopyBtn.Find(".copy-icon").Length(),
					"IPv6 Address copy icon not found",
				)
				assert.Equal(
					t,
					1,
					fullIPCopyBtn.Find(".copy-tooltip").Length(),
					"IPv6 Address copy tooltip not found",
				)

				// Verify copy button event listeners
				assert.Equal(
					t,
					1,
					doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
						return strings.Contains(
							s.Text(),
							`document.getElementById("copy-interface").addEventListener("click"`,
						)
					}).Length(),
					"Interface ID copy button event listener not found",
				)
				assert.Equal(
					t,
					1,
					doc.Find("script").FilterFunction(func(_ int, s *goquery.Selection) bool {
						return strings.Contains(
							s.Text(),
							`document.getElementById("copy-ip-full").addEventListener("click"`,
						)
					}).Length(),
					"IPv6 Address copy button event listener not found",
				)
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
				t.Helper()
				assert.Equal(
					t,
					"Invalid MAC address",
					doc.Find("p.error-message").Text(),
					"Incorrect error message",
				)
				assert.Equal(
					t,
					0,
					doc.Find("label[for='interface-id']").Length(),
					"Success fields should not be present",
				)
				assert.Equal(
					t,
					0,
					doc.Find("input#interface-id").Length(),
					"Success fields should not be present",
				)
				assert.Equal(
					t,
					0,
					doc.Find("label[for='ip-full']").Length(),
					"Success fields should not be present",
				)
				assert.Equal(
					t,
					0,
					doc.Find("input#ip-full").Length(),
					"Success fields should not be present",
				)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := renderToString(t, Result(tt.data))
			doc := parseHTML(t, html)
			tt.assertDoc(t, doc)
		})
	}
}
