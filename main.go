package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/LazyCode2/Koyo-site/config"
	"github.com/LazyCode2/Koyo-site/pages"
)

var (
	initFlag  = flag.Bool("init", false, "Initialize a new koyo-site project")
	buildFlag = flag.Bool("build", false, "Build the static site")
	serveFlag = flag.Bool("serve", false, "Serve the site locally")
	addFile   = flag.String("add", "", "Add a file")
)

func main() {
	// Parsing the flag for cli commands
	flag.Parse()

	switch {
	case *initFlag:
		initProject()
	case *buildFlag:
		buildSite()
	case *serveFlag:
		serveSite()
	case *addFile != "":
		addNewFile(*addFile)
	default:
		printHelp()
	}
}

func initProject() {
	dirs := []string{
		"content",
		"templates",
		"public",
	}

	fmt.Println("📁 Initializing koyo-site project...")

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ Failed to create %s: %v\n", dir, err)
			os.Exit(1)
		}
		fmt.Printf("✔ Created %s/\n", dir)
	}

	// Config file name
	configFile := "koyo.config.yaml"

configContent := `site:
  title: "My Koyo Site"
  author: "Your Name"
  bio: "Your Bio"

paths:
  content: "content"
  templates: "templates"
  output: "public"

server:
  port: ":8080"
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		fmt.Println("❌ Failed to create config file:", err)
		os.Exit(1)
	}

	fmt.Println("✔ Created koyo.config.yaml")
	fmt.Println("✨ koyo-site project initialized")
}

func addNewFile(filename string) {
	// Load config
	cfg, err := config.LoadConf()
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		os.Exit(1)
	}

	Content := "New post"
	//{content/filename.md}
	if err := os.WriteFile(cfg.Paths.Content+"/"+filename+".md", []byte(Content), 0644); err != nil {
		fmt.Println("❌ Failed to create" + filename + ".md")
		os.Exit(1)
	}
	fmt.Println("✔ Created " + filename + ".md")
}

func buildSite() {
	fmt.Println("⚙️  Building site...")

	// Load config
	cfg, err := config.LoadConf()
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		os.Exit(1)
	}

	blogsDir := filepath.Join(cfg.Paths.Output, "blogs")
	if err := os.MkdirAll(blogsDir, 0755); err != nil {
		fmt.Println("❌ Failed to create blogs directory:", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir(cfg.Paths.Content)
	if err != nil {
		fmt.Println("❌ Failed to read content directory:", err)
		os.Exit(1)
	}

	postTemplatePath := filepath.Join(cfg.Paths.Templates, "default.tmpl")
	if _, err := os.Stat(postTemplatePath); os.IsNotExist(err) {
		fmt.Println("❌ Template not found:", postTemplatePath)
		os.Exit(1)
	}

	// Build individual blog posts
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		if entry.Name() == "_index.md" {
			continue
		}

		contentPath := filepath.Join(cfg.Paths.Content, entry.Name())
		outputName := strings.TrimSuffix(entry.Name(), ".md") + ".html"
		outputPath := filepath.Join(blogsDir, outputName)

		fmt.Printf("📄 Building %s -> blogs/%s\n", entry.Name(), outputName)

		if err := pages.GeneratePage(contentPath, postTemplatePath, outputPath); err != nil {
			fmt.Printf("❌ Failed to generate %s: %v\n", entry.Name(), err)
			continue
		}
	}

	// Build index page
	indexTemplatePath := filepath.Join(cfg.Paths.Templates, "index.tmpl")
	if _, err := os.Stat(indexTemplatePath); os.IsNotExist(err) {
		fmt.Println("⚠️  index.tmpl not found, skipping index generation")
	} else {
		fmt.Println("📄 Building index.html")
		indexOutputPath := filepath.Join(cfg.Paths.Output, "index.html")

		if err := pages.GenerateIndexPage(
			cfg.Paths.Content,
			indexTemplatePath,
			indexOutputPath,
			cfg.Site.Title,
			cfg.Site.Author,
			cfg.Site.Bio,
		); err != nil {
			fmt.Printf("❌ Failed to generate index: %v\n", err)
		}
	}

	fmt.Println("✅ Site built successfully!")
}

func serveSite() {
	//Building site to avoid error while serving
	buildSite()
	// Load config
	cfg, err := config.LoadConf()
	if err != nil {
		fmt.Println("❌ Failed to load config:", err)
		os.Exit(1)
	}

	fs := http.FileServer(http.Dir(cfg.Paths.Output))
	http.Handle("/", fs)

	fmt.Println("🚀 Serving site at http://localhost:"+ cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, nil); err != nil {
		fmt.Println("❌ Failed to start server")
	}
}

func printHelp() {
	fmt.Println(`
koyo-site — a minimal static site generator

Usage:
  koyo-site [command]

Commands:
  -init           Initialize a new project
  -build          Build the site
  -serve          Serve locally
  -new <filename> Create markdown post at content
`)
	os.Exit(0)
}
