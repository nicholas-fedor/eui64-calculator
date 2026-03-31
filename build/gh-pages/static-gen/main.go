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
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// run performs the main logic of generating static HTML for the EUI-64 calculator.
// It renders the home page template, adapts it for static use by removing server-specific
// dependencies, adds WebAssembly scripts, formats the HTML for readability, and writes
// the output to dist/static/index.html. Returns an error if any step fails.
func run() error {
	// Create a buffer to hold the rendered HTML.
	var buf bytes.Buffer
	if err := ui.Home().Render(context.Background(), &buf); err != nil {
		return fmt.Errorf("failed to render home template: %w", err)
	}

	// Modify HTML for static site: remove HTMX, adjust paths, add WASM/JS scripts.
	htmlContent := buf.String()
	htmlContent = removeHTMXScript(htmlContent)
	htmlContent = replaceServerPaths(htmlContent)
	htmlContent = replaceLayoutScript(htmlContent)
	htmlContent = fixInputPatterns(htmlContent)

	// Format HTML for readability with newlines and indentation.
	formattedHTML, err := formatHTML(htmlContent)
	if err != nil {
		return fmt.Errorf("failed to format HTML: %w", err)
	}

	// Ensure output directory exists.
	outputDir := filepath.Join("dist", "static")
	if err := os.MkdirAll(outputDir, dirPerms); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
	}

	// Write formatted HTML to file.
	outputFile := filepath.Join(outputDir, "index.html")
	if err := os.WriteFile(outputFile, []byte(formattedHTML), filePerms); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputFile, err)
	}

	if _, err := fmt.Fprintln(
		os.Stdout,
		"Successfully generated static HTML at",
		outputFile,
	); err != nil {
		return fmt.Errorf("failed to write success message to stdout: %w", err)
	}

	return nil
}

// removeHTMXScript removes the HTMX script tag from the HTML, matching any version
// or integrity attribute to handle CDN updates.
func removeHTMXScript(htmlContent string) string {
	re := regexp.MustCompile(`<script\s+src="https://unpkg\.com/htmx\.org@[^"]+"[^>]*></script>`)

	return re.ReplaceAllString(htmlContent, "")
}

// replaceServerPaths updates server-specific paths and attributes to relative paths
// and static-compatible attributes for GitHub Pages deployment, unescaping inline
// script content to ensure valid JavaScript syntax.
func replaceServerPaths(htmlContent string) string {
	const minMatchCount = 2 // Minimum number of regex matches (full match + one submatch).
	// Remove all hx-* attributes from <form> tag.
	re := regexp.MustCompile(`(<form[^>]*?)(?:\s+hx-[^=\s]+="[^"]*")*(\s*[^>]*>)`)
	htmlContent = re.ReplaceAllString(htmlContent, "$1$2")

	// Adjust static asset paths for GitHub Pages.
	htmlContent = regexp.MustCompile(`/static/styles\.css`).ReplaceAllString(htmlContent, "./styles.css")
	htmlContent = regexp.MustCompile(`/static/favicon\.ico`).ReplaceAllString(htmlContent, "./favicon.ico")

	// Add preload hint for styles.css to prevent FOUC.
	htmlContent = regexp.MustCompile(`<head>`).
		ReplaceAllString(htmlContent, `<head><link rel="preload" href="./styles.css" as="style">`)

	// Unescape HTML entities in <script> tags to prevent issues with JavaScript execution.
	reScript := regexp.MustCompile(`<script\b[^>]*>(.*?)</script>`)

	return reScript.ReplaceAllStringFunc(htmlContent, func(script string) string {
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
func replaceLayoutScript(htmlContent string) string {
	// Match the entire inline script, including copyToClipboard and htmx:afterSwap, with dotall mode.
	formAttrRegex := regexp.MustCompile(
		`(?s)<script>\s*function copyToClipboard\(elementId, buttonId\)\s*{.*?document\.body\.addEventListener\('htmx:afterSwap'.*?</script>`,
	)
	if !formAttrRegex.MatchString(htmlContent) {
		fmt.Fprintf(
			os.Stderr,
			"Warning: Inline copyToClipboard script not found in HTML; WebAssembly scripts may not be included\n",
		)
	}

	return formAttrRegex.ReplaceAllString(htmlContent, `<script src="./wasm_exec.js"></script>
            <script src="./scripts.js"></script>`)
}

// fixInputPatterns fixes regex patterns in <input> elements to be compatible with
// modern browsers using the 'v' flag for pattern validation, avoiding character
// class ranges for the separator to prevent "invalid character in class" errors.
func fixInputPatterns(htmlContent string) string {
	// Match MAC address pattern: [0-9a-fA-F]{2}([-:][0-9a-fA-F]{2}){5}
	macPatternRegex := regexp.MustCompile(
		`pattern\s*=\s*"\[0\-9a\-fA\-F\]\{2\}\(\[\-:\]\[0\-9a\-fA\-F\]\{2\}\)\{5\}"`,
	)
	if macPatternRegex.MatchString(htmlContent) {
		fmt.Fprintf(os.Stderr, "Info: MAC address pattern matched in HTML; applying regex fix\n")

		return macPatternRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
			fmt.Fprintf(os.Stderr, "Debug: Matched MAC address pattern: %s\n", match)

			return `pattern="[0-9a-fA-F]{2}((-|:)[0-9a-fA-F]{2}){5}"`
		})
	}

	fmt.Fprintf(
		os.Stderr,
		"Warning: MAC address pattern not found in HTML; regex fix not applied\n",
	)

	return htmlContent
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
				fmt.Fprintf(&buf, " %s=\"%s\"", attr.Key, html.EscapeString(attr.Val))
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
				lines := strings.SplitSeq(text, "\n")
				for line := range lines {
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
