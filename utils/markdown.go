package utils

import "github.com/gomarkdown/markdown"

// Convert markdown to HTML
func MarkdownToHTML(content []byte) []byte {
	return markdown.ToHTML(content, nil, nil)
}
