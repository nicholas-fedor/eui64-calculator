// Package main generates static HTML for the EUI-64 calculator's GitHub Pages
// site. It renders UI templates, adapts them for client-side WebAssembly usage,
// formats the HTML for readability, and writes the result to a file for deployment.
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_main tests the main function's behavior by calling the run function directly.
// It verifies that run completes successfully, creates the output file, and generates
// valid HTML content with expected modifications for static hosting.
func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Successful run generates static HTML",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing output directory before test
			if err := os.RemoveAll("dist"); err != nil && !os.IsNotExist(err) {
				t.Fatalf("Failed to clean up dist directory: %v", err)
			}

			// Call run function directly (equivalent to main without os.Exit)
			err := run()
			require.NoError(t, err, "run() should complete without error")

			// Verify output file exists
			outputFile := filepath.Join("dist", "static", "index.html")
			assert.FileExists(t, outputFile, "Output file should be created")

			// Read and verify file content
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err, "Should be able to read output file")

			htmlContent := string(content)

			// Verify key transformations are present
			assert.Contains(t, htmlContent, "EUI-64 Calculator", "Should contain page title")
			assert.Contains(t, htmlContent, `<script src="./wasm_exec.js">`, "Should include WebAssembly runtime script")
			assert.Contains(t, htmlContent, `<script src="./scripts.js">`, "Should include application script")
			assert.Contains(t, htmlContent, `pattern="[0-9a-fA-F]{2}((-|:)[0-9a-fA-F]{2}){5}"`, "Should have fixed MAC address pattern")
			assert.NotContains(t, htmlContent, `src="https://unpkg.com/htmx.org`, "Should not contain HTMX CDN script")

			// Verify HTML is properly formatted (contains newlines and indentation)
			assert.Contains(t, htmlContent, "\n", "HTML should be formatted with newlines")
			assert.Contains(t, htmlContent, "  ", "HTML should be formatted with indentation")
		})
	}
}

// TestRemoveHTMXScript tests the removeHTMXScript function with various HTML inputs.
// It verifies that HTMX script tags are correctly removed, handling different versions and attributes.
func TestRemoveHTMXScript(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "HTMX script with version and integrity",
			html: `<html><head><script src="https://unpkg.com/htmx.org@1.9.10" integrity="sha384-..."></script></head><body></body></html>`,
			want: `<html><head></head><body></body></html>`,
		},
		{
			name: "HTMX script without integrity",
			html: `<html><head><script src="https://unpkg.com/htmx.org@2.0.0"></script></head><body></body></html>`,
			want: `<html><head></head><body></body></html>`,
		},
		{
			name: "No HTMX script",
			html: `<html><head><script src="https://example.com/script.js"></script></head><body></body></html>`,
			want: `<html><head><script src="https://example.com/script.js"></script></head><body></body></html>`,
		},
		{
			name: "Multiple scripts, only HTMX removed",
			html: `<html><head><script src="https://unpkg.com/htmx.org@1.9.10"></script><script src="https://example.com/script.js"></script></head><body></body></html>`,
			want: `<html><head><script src="https://example.com/script.js"></script></head><body></body></html>`,
		},
		{
			name: "Empty HTML",
			html: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeHTMXScript(tt.html)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestReplaceServerPaths tests the replaceServerPaths function with various HTML inputs.
// It verifies that server-specific paths and attributes are correctly adapted for static hosting.
func TestReplaceServerPaths(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "Replace static asset paths",
			html: `<html><head><link rel="stylesheet" href="/static/styles.css"><link rel="icon" href="/static/favicon.ico"></head><body></body></html>`,
			want: `<html><head><link rel="preload" href="./styles.css" as="style"><link rel="stylesheet" href="./styles.css"><link rel="icon" href="./favicon.ico"></head><body></body></html>`,
		},
		{
			name: "Remove hx-* attributes from form",
			html: `<form hx-post="/calculate" hx-target="#result" hx-swap="innerHTML"><input type="text"></form>`,
			want: `<form><input type="text"></form>`,
		},
		{
			name: "Unescape HTML entities in script",
			html: `<script>function test() { return ""hello" & <world>"; }</script>`,
			want: `<script>function test() { return ""hello" & <world>"; }</script>`,
		},
		{
			name: "Complex HTML with multiple modifications",
			html: `<html><head><link rel="stylesheet" href="/static/styles.css"><link rel="icon" href="/static/favicon.ico"></head><body><form hx-post="/calculate"><script>alert(""test"");</script></form></body></html>`,
			want: `<html><head><link rel="preload" href="./styles.css" as="style"><link rel="stylesheet" href="./styles.css"><link rel="icon" href="./favicon.ico"></head><body><form><script>alert(""test"");</script></form></body></html>`,
		},
		{
			name: "No modifications needed",
			html: `<html><head><link rel="stylesheet" href="./styles.css"></head><body><form><input type="text"></form></body></html>`,
			want: `<html><head><link rel="preload" href="./styles.css" as="style"><link rel="stylesheet" href="./styles.css"></head><body><form><input type="text"></form></body></html>`,
		},
		{
			name: "Empty HTML",
			html: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceServerPaths(tt.html)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestReplaceLayoutScript tests the replaceLayoutScript function with various HTML inputs.
// It verifies that the inline JavaScript is correctly replaced with WebAssembly script includes.
func TestReplaceLayoutScript(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "Replace inline script with WebAssembly includes",
			html: `<script>
function copyToClipboard(elementId, buttonId) {
    // copy logic
}
document.body.addEventListener('htmx:afterSwap', function() {
    // htmx logic
});
</script>`,
			want: `<script src="./wasm_exec.js"></script>
            <script src="./scripts.js"></script>`,
		},
		{
			name: "No matching script",
			html: `<script>console.log("other script");</script>`,
			want: `<script>console.log("other script");</script>`,
		},
		{
			name: "Empty HTML",
			html: "",
			want: "",
		},
		{
			name: "Script without htmx event listener",
			html: `<script>function copyToClipboard() { /* logic */ }</script>`,
			want: `<script>function copyToClipboard() { /* logic */ }</script>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceLayoutScript(tt.html)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestFixInputPatterns tests the fixInputPatterns function with various HTML inputs.
// It verifies that MAC address input patterns are correctly fixed for browser compatibility.
func TestFixInputPatterns(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "Fix MAC address pattern",
			html: `<input type="text" pattern="[0-9a-fA-F]{2}([-:][0-9a-fA-F]{2}){5}">`,
			want: `<input type="text" pattern="[0-9a-fA-F]{2}((-|:)[0-9a-fA-F]{2}){5}">`,
		},
		{
			name: "No pattern to fix",
			html: `<input type="text" pattern="[0-9]+">`,
			want: `<input type="text" pattern="[0-9]+">`,
		},
		{
			name: "Empty HTML",
			html: "",
			want: "",
		},
		{
			name: "Multiple inputs, only MAC pattern fixed",
			html: `<input type="text" pattern="[0-9a-fA-F]{2}([-:][0-9a-fA-F]{2}){5}"><input type="text" pattern="[0-9]+">`,
			want: `<input type="text" pattern="[0-9a-fA-F]{2}((-|:)[0-9a-fA-F]{2}){5}"><input type="text" pattern="[0-9]+">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fixInputPatterns(tt.html)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestFormatHTML tests the formatHTML function with various HTML inputs.
// It verifies that HTML is correctly parsed and formatted with proper indentation and newlines.
func TestFormatHTML(t *testing.T) {
	tests := []struct {
		name    string
		htmlStr string
		want    string
		wantErr bool
	}{
		{
			name:    "Simple HTML formatting",
			htmlStr: `<html><head><title>Test</title></head><body><h1>Hello</h1></body></html>`,
			want:    "<html>\n  <head>\n    <title>\n      Test\n    </title>\n  </head>\n  <body>\n    <h1>\n      Hello\n    </h1>\n  </body>\n</html>\n",
			wantErr: false,
		},
		{
			name:    "HTML with void elements",
			htmlStr: `<html><head><meta charset="utf-8"><link rel="stylesheet" href="style.css"></head><body><img src="image.jpg" alt="test"><br></body></html>`,
			want:    "<html>\n  <head>\n    <meta charset=\"utf-8\"/>\n    <link rel=\"stylesheet\" href=\"style.css\"/>\n  </head>\n  <body>\n    <img src=\"image.jpg\" alt=\"test\"/>\n    <br/>\n  </body>\n</html>\n",
			wantErr: false,
		},
		{
			name:    "HTML with script content",
			htmlStr: `<html><head><script>console.log("test");</script></head><body></body></html>`,
			want:    "<html>\n  <head>\n    <script>\n      console.log(\"test\");\n    </script>\n  </head>\n  <body>\n  </body>\n</html>\n",
			wantErr: false,
		},
		{
			name:    "HTML with comments",
			htmlStr: `<html><!-- comment --><body></body></html>`,
			want:    "<html>\n  <!--comment-->\n  <head>\n  </head>\n  <body>\n  </body>\n</html>\n",
			wantErr: false,
		},
		{
			name:    "Invalid HTML",
			htmlStr: `<html><head><title>Test</title></head><body><h1>Hello</h1>`,
			want:    "<html>\n  <head>\n    <title>\n      Test\n    </title>\n  </head>\n  <body>\n    <h1>\n      Hello\n    </h1>\n  </body>\n</html>\n",
			wantErr: false,
		},
		{
			name:    "Empty string",
			htmlStr: "",
			want:    "<html>\n  <head>\n  </head>\n  <body>\n  </body>\n</html>\n",
			wantErr: false,
		},
		{
			name:    "HTML with nested elements",
			htmlStr: `<div><p><span>text</span></p></div>`,
			want:    "<html>\n  <head>\n  </head>\n  <body>\n    <div>\n      <p>\n        <span>\n          text\n        </span>\n      </p>\n    </div>\n  </body>\n</html>\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatHTML(tt.htmlStr)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
