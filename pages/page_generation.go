package pages

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	parse "github.com/LazyCode2/Koyo-site/parser"
	"github.com/gomarkdown/markdown"
)

// Page is a struct for single page with frontmatter and content
type Page struct {
	Title       string
	Description string
	Author      string
	Date        string
	Content     template.HTML
	Meta        map[string]interface{} // additional frontmatter fields
}

func BuildPage(contentPath string) *Page {
	content, _, err := parse.GetContent(contentPath)

	if err != nil {
		log.Fatal("cannot read markdown:", err)
	}

	frontmatter, bodyContent := parse.ParseFrontmatter(content)
	htmlBody := markdown.ToHTML(bodyContent, nil, nil)

	page := &Page{
		Content: template.HTML(htmlBody),
		Meta:    frontmatter,
	}

	// Extract common fields from frontmatter
	if frontmatter != nil {
		if title, ok := frontmatter["title"].(string); ok {
			page.Title = title
		}
		if desc, ok := frontmatter["description"].(string); ok {
			page.Description = desc
		}
		if author, ok := frontmatter["author"].(string); ok {
			page.Author = author
		}
		if date, ok := frontmatter["date"].(string); ok {
			page.Date = date
		}
	}

	return page
}

func RenderPage(page *Page, templatePath string) ([]byte, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, page); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

func GeneratePage(contentPath, templatePath, outputPath string) error {
	page := BuildPage(contentPath)

	html, err := RenderPage(page, templatePath)
	if err != nil {
		return err
	}
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, html, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
