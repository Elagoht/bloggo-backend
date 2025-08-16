package readtime

import (
	"regexp"
	"strings"
)

// EstimateReadTime calculates estimated reading time in minutes for given content
// Based on average reading speed of 200 words per minute
func EstimateReadTime(content string) int {
	if content == "" {
		return 0
	}

	// Remove HTML tags if any
	re := regexp.MustCompile(`<[^>]*>`)
	cleanContent := re.ReplaceAllString(content, "")

	// Remove extra whitespace and split into words
	cleanContent = strings.TrimSpace(cleanContent)
	words := strings.Fields(cleanContent)
	wordCount := len(words)

	if wordCount == 0 {
		return 0
	}

	// Calculate reading time (200 words per minute)
	const averageWPM = 200
	readTimeMinutes := max(
		// Round up
		// Minimum read time is 1 minute
		(wordCount+averageWPM-1)/averageWPM,
		1,
	)

	return readTimeMinutes
}
