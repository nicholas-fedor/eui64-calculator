// Package main generates static HTML for the EUI-64 calculator's GitHub Pages
// site. It renders UI templates, adapts them for client-side WebAssembly usage,
// formats the HTML for readability, and writes the result to a file for deployment.
package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/nicholas-fedor/eui64-calculator/internal/ui"
)

// dirPerms defines the permission bits for creating directories, ensuring
// read/write/execute for owner and group, and read/execute for others.
const dirPerms = 0o755

// filePerms defines the permission bits for writing files, allowing read access
// for all to suit static hosting requirements.
const filePerms = 0o644

// main renders the EUI-64 calculator's home page template, adapts it for static use
// by removing server-specific dependencies, adds WebAssembly scripts, formats the
// HTML for readability, and writes the output to dist/static/index.html.
func main() {
	// Create a buffer to hold the rendered HTML.
	var buf bytes.Buffer
	if err := ui.Home().Render(context.Background(), &buf); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to render home template: %v\n", err)
		os.Exit(1)
	}

	// Modify HTML for static site: remove HTMX, adjust paths, add WASM/JS scripts.
	html := buf.String()
	html = removeHTMXScript(html)
	html = replaceServerPaths(html)
	html = replaceLayoutScript(html)
	html = fixInputPatterns(html)

	// Format HTML for readability with newlines and indentation.
	formattedHTML, err := formatHTML(html)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to format HTML: %v\n", err)
		os.Exit(1)
	}

	// Ensure output directory exists.
	outputDir := filepath.Join("dist", "static")
	if err := os.MkdirAll(outputDir, dirPerms); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	// Write formatted HTML to file.
	outputFile := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputFile, []byte(formattedHTML), filePerms); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", outputFile, err)
		os.Exit(1)
	}

	// Verify required static assets exist and check main.wasm integrity.
	assets := []string{"styles.css", "favicon.ico", "wasm_exec.js", "scripts.js", "main.wasm"}
	for _, asset := range assets {
		assetPath := filepath.Join(outputDir, asset)
		if _, err := os.Stat(assetPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Required asset %s not found in %s\n", asset, outputDir)
			os.Exit(1)
		}

		if asset == "main.wasm" {
			data, err := os.ReadFile(assetPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to read %s: %v\n", assetPath, err)
				os.Exit(1)
			}
			// Check for WebAssembly magic number (0x00 0x61 0x73 0x6D).
			if len(data) < 4 || data[0] != 0x00 || data[1] != 0x61 || data[2] != 0x73 ||
				data[3] != 0x6D {
				fmt.Fprintf(
					os.Stderr,
					"Error: %s is not a valid WebAssembly binary (first 4 bytes: %x %x %x %x, expected: 00 61 73 6d)\n",
					assetPath,
					data[0],
					data[1],
					data[2],
					data[3],
				)
				os.Exit(1)
			}
		}
	}

	fmt.Fprintln(os.Stdout, "Successfully generated static HTML at", outputFile)
}

// removeHTMXScript removes the HTMX script tag from the HTML, matching any version
// or integrity attribute to handle CDN updates.
func removeHTMXScript(html string) string {
	re := regexp.MustCompile(`<script\s+src="https://unpkg\.com/htmx\.org@[^"]+"[^>]*></script>`)

	return re.ReplaceAllString(html, "")
}

// replaceServerPaths updates server-specific paths and attributes to relative paths
// and static-compatible attributes for GitHub Pages deployment, unescaping inline
// script content to ensure valid JavaScript syntax.
func replaceServerPaths(html string) string {
	const minMatchCount = 2 // Minimum number of regex matches (full match + one submatch).
	// Remove all hx-* attributes from <form> tag.
	re := regexp.MustCompile(`(<form[^>]*?)(?:\s+hx-[^=\s]+="[^"]*")*(\s*[^>]*>)`)
	html = re.ReplaceAllString(html, "$1$2")

	// Adjust static asset paths for GitHub Pages.
	html = regexp.MustCompile(`/static/styles\.css`).ReplaceAllString(html, "./styles.css")
	html = regexp.MustCompile(`/static/favicon\.ico`).ReplaceAllString(html, "./favicon.ico")

	// Add preload hint for styles.css to prevent FOUC.
	html = regexp.MustCompile(`<head>`).
		ReplaceAllString(html, `<head><link rel="preload" href="./styles.css" as="style">`)

	// Unescape HTML entities in <script> tags to prevent issues with JavaScript execution.
	reScript := regexp.MustCompile(`<script\b[^>]*>(.*?)</script>`)

	return reScript.ReplaceAllStringFunc(html, func(script string) string {
		// Extract the script content (between <script> and </script>).
		reContent := regexp.MustCompile(`<script\b[^>]*>(.*?)</script>`)

		matches := reContent.FindStringSubmatch(script)
		if len(matches) < minMatchCount {
			return script // Return unchanged if no content found.
		}

		content := matches[1]
		// Unescape common HTML entities.
		content = strings.ReplaceAll(content, "&quot;", "\"")
		content = strings.ReplaceAll(content, "&amp;", "&")
		content = strings.ReplaceAll(content, "&lt;", "<")
		content = strings.ReplaceAll(content, "&gt;", ">")

		return fmt.Sprintf("<script>%s</script>", content)
	})
}

// replaceLayoutScript replaces the layout's inline JavaScript with includes for
// WebAssembly runtime and application scripts, enabling client-side functionality.
func replaceLayoutScript(html string) string {
	// Match the entire inline script, including copyToClipboard and htmx:afterSwap, with dotall mode.
	formAttrRegex := regexp.MustCompile(
		`(?s)<script>\s*function copyToClipboard\(elementId, buttonId\)\s*{.*?document\.body\.addEventListener\('htmx:afterSwap'.*?</script>`,
	)
	if !formAttrRegex.MatchString(html) {
		fmt.Fprintf(
			os.Stderr,
			"Warning: Inline copyToClipboard script not found in HTML; WebAssembly scripts may not be included\n",
		)
	}

	return formAttrRegex.ReplaceAllString(html, `<script src="./wasm_exec.js"></script>
            <script src="./scripts.js"></script>`)
}

// fixInputPatterns fixes regex patterns in <input> elements to be compatible with
// modern browsers using the 'v' flag for pattern validation, avoiding character
// class ranges for the separator to prevent "invalid character in class" errors.
func fixInputPatterns(html string) string {
	// Match MAC address pattern: [0-9a-fA-F]{2}([-:][0-9a-fA-F]{2}){5}
	macPatternRegex := regexp.MustCompile(
		`pattern\s*=\s*"\[0\-9a\-fA\-F\]\{2\}\(\[\-:\]\[0\-9a\-fA\-F\]\{2\}\)\{5\}"`,
	)
	if macPatternRegex.MatchString(html) {
		fmt.Fprintf(os.Stderr, "Info: MAC address pattern matched in HTML; applying regex fix\n")

		return macPatternRegex.ReplaceAllStringFunc(html, func(match string) string {
			fmt.Fprintf(os.Stderr, "Debug: Matched MAC address pattern: %s\n", match)

			return `pattern="[0-9a-fA-F]{2}((-|:)[0-9a-fA-F]{2}){5}"`
		})
	}

	fmt.Fprintf(
		os.Stderr,
		"Warning: MAC address pattern not found in HTML; regex fix not applied\n",
	)

	return html
}

// formatHTML parses the input HTML and returns a formatted version with newlines
// and 2-space indentation for readability, using the html package to ensure
// consistent tag separation.
func formatHTML(htmlStr string) (string, error) {
	var buf bytes.Buffer

	node, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// voidElements defines HTML elements that are self-closing and do not require
	// a closing tag, per the HTML specification.
	voidElements := map[string]bool{
		"area": true, "base": true, "br": true, "col": true, "embed": true,
		"hr": true, "img": true, "input": true, "link": true, "meta": true,
		"param": true, "source": true, "track": true, "wbr": true,
	}

	// renderNode recursively renders an HTML node with 2-space indentation and
	// newlines after tags and text content.
	var renderNode func(node *html.Node, depth int)

	renderNode = func(node *html.Node, depth int) {
		indent := strings.Repeat("  ", depth)

		switch node.Type {
		case html.DoctypeNode:
			buf.WriteString("<!DOCTYPE html>\n")
		case html.ElementNode:
			// Start tag.
			buf.WriteString(indent + "<" + node.Data)

			for _, attr := range node.Attr {
				// Escape attribute values to prevent XSS.
				buf.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Key, html.EscapeString(attr.Val)))
			}

			if voidElements[node.Data] {
				buf.WriteString("/>\n")

				return
			}

			buf.WriteString(">\n")

			// Handle special case for script tags to preserve content verbatim.
			if node.Data == "script" {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						// Write script content without escaping to preserve JavaScript.
						buf.WriteString(indent + "  " + strings.TrimSpace(c.Data) + "\n")
					} else {
						renderNode(c, depth+1)
					}
				}
			} else {
				// Render child nodes normally.
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					renderNode(c, depth+1)
				}
			}

			// End tag.
			buf.WriteString(indent + "</" + node.Data + ">\n")
		case html.TextNode:
			text := strings.TrimSpace(node.Data)
			if text != "" {
				lines := strings.Split(text, "\n")
				for _, line := range lines {
					if line = strings.TrimSpace(line); line != "" {
						// Escape text content to prevent XSS, unless it's within a script tag.
						if node.Parent != nil && node.Parent.Data != "script" {
							buf.WriteString(indent + html.EscapeString(line) + "\n")
						} else {
							buf.WriteString(indent + line + "\n")
						}
					}
				}
			}
		case html.CommentNode:
			if comment := strings.TrimSpace(node.Data); comment != "" {
				buf.WriteString(indent + "<!--" + html.EscapeString(comment) + "-->\n")
			}
		case html.DocumentNode:
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				renderNode(c, depth)
			}
		case html.ErrorNode, html.RawNode:
			// Skip error and raw nodes, as they are not expected in valid HTML.
		}
	}

	renderNode(node, 0)

	return strings.TrimSpace(buf.String()) + "\n", nil
}
