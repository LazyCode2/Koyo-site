package pages

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	parse "github.com/LazyCode2/Koyo-site/parser"
	"github.com/gomarkdown/markdown"
)

// Page represents a single rendered page
type Page struct {
	Title       string
	Description string
	Author      string
	Date        string
	Content     template.HTML
	Meta        Meta
}

type Meta struct {
	Tags      []string
	SiteTitle string
}

func BuildPage(contentPath string) (*Page, error) {
	content, _, err := parse.GetContent(contentPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read markdown: %w", err)
	}

	frontmatter, bodyContent := parse.ParseFrontmatter(content)
	htmlBody := markdown.ToHTML(bodyContent, nil, nil)

	page := &Page{
		Content: template.HTML(htmlBody),
	}

	// Frontmatter fields
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

		// Tags (string or list)
		if rawTags, ok := frontmatter["tags"]; ok {
			page.Meta.Tags = parse.ParseTagsRaw(rawTags)
		}
	}

	return page, nil
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
	page, err := BuildPage(contentPath)
	if err != nil {
		return err
	}

	html, err := RenderPage(page, templatePath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, html, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
