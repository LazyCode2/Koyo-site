package pages

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/LazyCode2/Koyo-site/utils"
	"github.com/gomarkdown/markdown"
)

// PostMeta is a struct fpr blog post listing
type PostMeta struct {
	Title       string
	Description string
	Author      string
	Date        string
	URL         string
	Filename    string
}

// IndexPage is the index page with posts listing
type IndexPage struct {
	SiteTitle     string
	SiteAuthor    string
	SiteAuthorBio string
	Content       template.HTML
	Posts         []PostMeta
}

var logger = utils.NewLogger()

func CollectPosts(contentDir string) ([]PostMeta, error) {
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, err
	}

	var posts []PostMeta

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		// Skip _index.md
		if entry.Name() == "_index.md" {
			continue
		}

		contentPath := filepath.Join(contentDir, entry.Name())
		content, err := os.ReadFile(contentPath)
		if err != nil {
			logger.Warn("Skipping %s: failed to read file: %v", entry.Name(), err)
			continue
		}

		frontmatter, _ := ParseFrontmatter(content)

		post := PostMeta{
			Filename: entry.Name(),
			URL:      "/blogs/" + strings.TrimSuffix(entry.Name(), ".md") + ".html",
		}

		if frontmatter != nil {
			if title, ok := frontmatter["title"].(string); ok {
				post.Title = title
			}
			if desc, ok := frontmatter["description"].(string); ok {
				post.Description = desc
			}
			if author, ok := frontmatter["author"].(string); ok {
				post.Author = author
			}
			if date, ok := frontmatter["date"].(string); ok {
				post.Date = date
			}
		}

		// Use filename as fallback title
		if post.Title == "" {
			post.Title = strings.TrimSuffix(entry.Name(), ".md")
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		dateI, errI := time.Parse("2006-01-02", posts[i].Date)
		dateJ, errJ := time.Parse("2006-01-02", posts[j].Date)

		if errI != nil || errJ != nil {
			return posts[i].Date > posts[j].Date
		}

		return dateI.After(dateJ)
	})

	return posts, nil
}

func GenerateIndexPage(contentDir, templatePath, outputPath, siteTitle, siteAuthor, siteAuthorBio string) error {
	// Read _index.md
	indexPath := filepath.Join(contentDir, "_index.md")
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read _index.md: %w", err)
	}

	frontmatter, bodyContent := ParseFrontmatter(indexContent)
	htmlBody := MarkdownToHTML(bodyContent)

	posts, err := CollectPosts(contentDir)
	if err != nil {
		return fmt.Errorf("failed to collect posts: %w", err)
	}

	// Build index page data
	indexPage := &IndexPage{
		SiteTitle:     siteTitle,
		SiteAuthor:    siteAuthor,
		SiteAuthorBio: siteAuthorBio,
		Content:       template.HTML(htmlBody),
		Posts:         posts,
	}

	// Override with frontmatter if present
	if frontmatter != nil {
		if title, ok := frontmatter["title"].(string); ok {
			indexPage.SiteTitle = title
		}
		if author, ok := frontmatter["author"].(string); ok {
			indexPage.SiteAuthor = author
		}
		if authorbio, ok := frontmatter["bio"].(string); ok {
			indexPage.SiteAuthorBio = authorbio
		}
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Error("Template parsing failed: %v", err)
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, indexPage); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// Helper function to convert markdown to HTML
func MarkdownToHTML(content []byte) []byte {
	return markdown.ToHTML(content, nil, nil)
}
