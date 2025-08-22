package slugify

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	slug := strings.ToLower(text)

	// Replace special characters with words
	replacements := map[string]string{
		"+": "plus",
		"#": "sharp",
		"&": "and",
		"@": "at",
	}

	for char, word := range replacements {
		slug = strings.ReplaceAll(slug, char, "-"+word+"-")
	}

	// Replace any remaining non-alphanumeric characters with dashes
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")

	// Clean up multiple dashes and trim
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}
