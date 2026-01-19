package common

import (
	"strings"
)

// CleanText removes extra whitespace and trims a string
func CleanText(text string) string {
	// Remove newlines and extra spaces
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Trim and collapse multiple spaces into one
	return strings.TrimSpace(strings.Join(strings.Fields(text), " "))
}

// CleanImageURL removes query parameters from image URLs
// Example: "image.jpg?w=300&resize=150" -> "image.jpg"
func CleanImageURL(url string) string {
	if strings.Contains(url, "?") {
		return strings.Split(url, "?")[0]
	}
	return url
}

// RemoveLazyPlaceholder checks if image URL is a lazy load placeholder
func RemoveLazyPlaceholder(url string) string {
	if strings.Contains(url, "lazy.jpg") || url == "" {
		return ""
	}
	return url
}

// ExtractFirstNonEmpty returns the first non-empty string from a list
func ExtractFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
