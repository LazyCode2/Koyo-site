package parse

import (
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseFrontmatter(content []byte) (map[string]interface{}, []byte) {
	strContent := string(content)

	if !strings.HasPrefix(strContent, "---") {
		return nil, content
	}

	// closing ---
	parts := strings.SplitN(strContent[3:], "---", 2)
	if len(parts) != 2 {
		return nil, content
	}

	// Parse YAML frontmatter
	frontmatter := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(parts[0]), &frontmatter); err != nil {
		log.Printf("Warning: failed to parse frontmatter: %v", err)
		return nil, content
	}

	return frontmatter, []byte(strings.TrimSpace(parts[1]))
}
