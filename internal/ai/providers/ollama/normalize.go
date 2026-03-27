package ollama

import (
	"regexp"
	"strings"
)

var thinkRe = regexp.MustCompile(`(?s)<think>.*?</think>`)

func normalize(s string) string {
	return strings.TrimSpace(thinkRe.ReplaceAllString(s, ""))
}
