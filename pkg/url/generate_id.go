package url

import (
	"log"
	"regexp"
	"strings"
)

// given a string of text, returns a url-safe string identifier.
func GenerateID(text string) string {
	maxLen := 32
	// all alphanumeric chars including space and hyphen
	r, err := regexp.CompilePOSIX("[^a-zA-Z0-9 -]+")
	if err != nil {
		log.Fatal("unable to parse id regular expression. Exiting.")
	}

	// 1. remove additional text lines
	firstLine := strings.Split(text, "\n")[0]
	// 2. remove symbols
	sanitizedLine := r.ReplaceAllString(firstLine, "")
	// 3. trim prefix and suffix spaces
	trimmedLine := strings.TrimSpace(sanitizedLine)
	// 4. convert spaces to hyphens
	dashedLine := strings.ReplaceAll(trimmedLine, " ", "-")
	// 5. convert to lowercase
	result := strings.ToLower(dashedLine)

	// 6. trim length and remove suffix hyphen if necessary
	if len(result) > maxLen {
		result = result[:maxLen]
	}
	return strings.TrimSuffix(result, "-")
}
