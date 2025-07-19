package slugify

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	slug := strings.ToLower(text)
	slug = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
