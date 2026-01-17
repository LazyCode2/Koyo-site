package main

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/LazyCode2/Koyo-site/config"
	"github.com/LazyCode2/Koyo-site/pages"
	"github.com/LazyCode2/Koyo-site/utils"
)

var (
	initFlag  = flag.Bool("init", false, "Initialize a new koyo-site project")
	buildFlag = flag.Bool("build", false, "Build the static site")
	serveFlag = flag.Bool("serve", false, "Serve the site locally")
	addFile   = flag.String("add", "", "Add a file")
)

var logger = utils.NewLogger()

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

	logger.Info("Initializing koyo-site project...")

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Fatal("❌ Failed to create %s: %v\n", dir, err)
		}
		logger.Info("✔ Created %s/\n", dir)
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
		logger.Fatal("❌ Failed to create config file:%v\n", err)
	}

	logger.Info("✔ Created koyo.config.yaml")
	logger.Info("koyo-site project initialized")
}

func addNewFile(filename string) error {
	// Load config
	cfg, err := config.LoadConf()
	if err != nil {
		logger.Fatal("❌ Failed to load config:%v\n", err)
	}

	content := "New post"
	filePath := filepath.Join(cfg.Paths.Content, filename+".md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}
	logger.Info("Created file: %s", filePath)
	return nil
}

func buildSite() error {
	logger.Info("Building site...")

	cfg, err := config.LoadConf()
	if err != nil {
		return err
	}

	blogsDir := filepath.Join(cfg.Paths.Output, "blogs")
	if err := os.MkdirAll(blogsDir, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(cfg.Paths.Content)
	if err != nil {
		return err
	}

	postTemplatePath := filepath.Join(cfg.Paths.Templates, "default.tmpl")
	if _, err := os.Stat(postTemplatePath); os.IsNotExist(err) {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == "_index.md" {
			continue
		}

		contentPath := filepath.Join(cfg.Paths.Content, entry.Name())
		outputName := strings.TrimSuffix(entry.Name(), ".md") + ".html"
		outputPath := filepath.Join(blogsDir, outputName)

		logger.Info("Building post: %s -> blogs/%s", entry.Name(), outputName)
		if err := pages.GeneratePage(contentPath, postTemplatePath, outputPath); err != nil {
			logger.Warn("Failed to generate %s: %v", entry.Name(), err)
		}
	}

	indexTemplatePath := filepath.Join(cfg.Paths.Templates, "index.tmpl")
	if _, err := os.Stat(indexTemplatePath); os.IsNotExist(err) {
		logger.Warn("index.tmpl not found, skipping index generation")
	} else {
		indexOutputPath := filepath.Join(cfg.Paths.Output, "index.html")
		logger.Info("Building index.html")

		if err := pages.GenerateIndexPage(cfg.Paths.Content, indexTemplatePath, indexOutputPath,
			cfg.Site.Title, cfg.Site.Author, cfg.Site.Bio); err != nil {
			logger.Warn("Failed to generate index: %v", err)
		}
	}

	logger.Info("Site built successfully")
	return nil
}

func serveSite() error {
	// Build first
	if err := buildSite(); err != nil {
		return err
	}

	cfg, err := config.LoadConf()
	if err != nil {
		return err
	}

	fs := http.FileServer(http.Dir(cfg.Paths.Output))
	http.Handle("/", fs)

	logger.Info("Serving site at http://localhost%s", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, nil); err != nil {
		return err
	}
	return nil
}

func printHelp() {
	logger.Info(`
koyo-site — a minimal static site generator

Usage:
  koyo-site [command]

Commands:
  -init           Initialize a new project
  -build          Build the site
  -serve          Serve locally
  -add <filename> Create markdown post in content/
`)
	os.Exit(0)
}
