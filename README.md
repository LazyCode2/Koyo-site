# Koyo-site

A minimalistic static site generator written in Go. Simple, fast, and focused on blogging.

## Screenshot

<img width="800" height="400" alt="Screenshot Welcome to My Blog" src="./assets/Welcome Blog.png" />

## Installation

### Install from source

```bash
git clone https://github.com/LazyCode2/Koyo-site.git
cd Koyo-site
go install
``` 

### Install directly

```bash
go install github.com/LazyCode2/Koyo-site@latest
```

### Build binary

```bash
go build -o Koyo-site main.go
```

## Quick Start

### 1. Initialize a new project

```bash
Koyo-site -init
```

This creates:
```
.
├── content/          # Your markdown files
├── templates/        # HTML templates
├── public/           # Generated site (output)
└── koyo.config.yaml  # Configuration file
```

### 2. Create your homepage

Create `content/_index.md`:

```markdown
---
title: "Welcome to My Blog"
---

# Hello! 👋

I'm a developer writing about tech and life.
```

### 3. Write your first post

Create `content/my-first-post.md`:

```markdown
---
title: "My First Post"
description: "Getting started with Koyo-site"
author: "Your Name"
date: "2024-12-15"
tags: 
  - Welcome
---

# My First Post

This is my first blog post using **Koyo-site**!

## Why I chose Koyo

- Simple and fast
- Minimal configuration
- Easy to customize
```

### 4. Create templates

**`templates/index.tmpl`** (homepage with post listings):

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.SiteTitle}}</title>
</head>
<body>
    <header>
        <h1>{{.SiteTitle}}</h1>
        {{if .SiteAuthor}}<p>By {{.SiteAuthor}}</p>{{end}}
    </header>

    <div class="intro">
        {{.Content}}
    </div>

    <section>
        <h2>Posts</h2>
        <ul>
            {{range .Posts}}
            <li>
                <a href="{{.URL}}">{{.Title}}</a>
                {{if .Date}} - {{.Date}}{{end}}
                {{if .Description}}<p>{{.Description}}</p>{{end}}
            </li>
            {{end}}
        </ul>
    </section>
</body>
</html>
```

**`templates/default.tmpl`** (individual blog posts):

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    {{if .Description}}<meta name="description" content="{{.Description}}">{{end}}
</head>
<body>
    <main>
        {{if .Title}}<h1>{{.Title}}</h1>{{end}}
        {{if .Date}}<time>{{.Date}}</time>{{end}}
        {{if .Author}}<p>By {{.Author}}</p>{{end}}
        
        <article>
            {{.Content}}
        </article>
    </main>
</body>
</html>
```

### 5. Build your site

```bash
Koyo-site -build
```

Output structure:
```
public/
├── index.html
└── blogs/
    └── my-first-post.html
```

### 6. Serve locally

```bash
Koyo-site -serve
```

## Configuration

Edit `koyo.config.yaml`:

```yaml
site:
  title: "My Koyo Site"
  author: "Your Name"
  bio: "Your Bio"

paths:
  content: "content"
  templates: "templates"
  output: "public"
```

## Template Variables

### Index Page (`index.tmpl`)

- `{{.SiteTitle}}` - Site title from config
- `{{.SiteAuthor}}` - Site author from config
- `{{.SiteAuthorBio}}` - Site author Bio from config
- `{{.Content}}` - HTML content from `_index.md`
- `{{.Posts}}` - Array of all posts

### Post Object (in `{{range .Posts}}`)

- `{{.Title}}` - Post title
- `{{.Description}}` - Post description
- `{{.Author}}` - Post author
- `{{.Date}}` - Post date
- `{{.URL}}` - Post URL (e.g., `/blogs/post-name.html`)

### Blog Post Page (`default.tmpl`)

- `{{.Title}}` - Post title
- `{{.Description}}` - Post description
- `{{.Author}}` - Post author
- `{{.Date}}` - Post date
- `{{.Content}}` - HTML content
- `{{.Meta}}` - Map of all frontmatter fields
- `{{.Meta.Tags}}` - Map all tags

## Frontmatter

All fields are optional:

```yaml
---
title: "Post Title"
description: "A brief description"
author: "Author Name"
date: "2024-12-15"
# Add custom fields
tags: 
  - go 
  - tech
draft: false
---
```

Access custom fields in templates via `{{.Meta}}`:

```html
{{if .Meta.Tags}}
  Tags: {{range .Meta.Tags}}{{.}}{{end}}
{{end}}
```

## Project Structure

```
your-project/
├── content/
│   ├── _index.md           # Homepage content (special file)
│   ├── first-post.md       # Blog posts
│   └── second-post.md
├── templates/
│   ├── index.tmpl          # Homepage template
│   └── default.tmpl        # Blog post template
├── public/                 # Generated output
│   ├── index.html
│   └── blogs/
│       ├── first-post.html
│       └── second-post.html
└── koyo.config.yaml        # Configuration
```

## Commands

```bash
Koyo-site -init     # Initialize new project
Koyo-site -build    # Build static site
Koyo-site -serve    # Serve locally
Koyo-site -add <filename> # Add post in content
```

## Dependencies

- [gomarkdown/markdown](https://github.com/gomarkdown/markdown) - Markdown parser
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) - YAML parser

## Contributing

Contributions are welcome! This is a minimalist tool, so new features should align with the goal of simplicity.

## License

MIT License - feel free to use and modify!

## Roadmap

- [ ] Hot reload during development
- [ ] Custom template functions
- [ ] RSS feed generation
- [ ] Sitemap generation
- [ ] Draft post support
- [ x ] Tag rendering

---

Built with ❤️ and Go
