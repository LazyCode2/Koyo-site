package parse

import "strings"

func ParseTags(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	tags := make([]string, 0, len(parts))

	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

func ParseTagsRaw(v interface{}) []string {
	switch t := v.(type) {
	case string:
		return ParseTags(t)
	case []string:
		return t
	case []interface{}:
		var tags []string
		for _, item := range t {
			if s, ok := item.(string); ok {
				tags = append(tags, s)
			}
		}
		return tags
	default:
		return nil
	}
}
