package branchname

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	dashDash   = regexp.MustCompile(`-{2,}`)
	allowedRe  = regexp.MustCompile(`[^a-z0-9-]`)
	whitespace = regexp.MustCompile(`\s+`)
)

func SlugifyDescription(desc string) string {
	s := strings.ToLower(desc)
	s = strings.ReplaceAll(s, " - ", " ")
	s = whitespace.ReplaceAllString(s, "-")
	s = allowedRe.ReplaceAllString(s, "")
	s = dashDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func normalizeType(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}

// BranchType returns the prefix before the first "/" in branch.
// Returns an error if no "/" is present.
func BranchType(branch string) (string, error) {
	idx := strings.Index(branch, "/")
	if idx < 0 {
		return "", fmt.Errorf("branchname: no '/' found in branch %q", branch)
	}
	return branch[:idx], nil
}

func BuildName(taskType, issue, description string) string {
	slug := SlugifyDescription(description)

	parts := []string{}
	if issue != "" {
		parts = append(parts, issue)
	}
	if slug != "" {
		parts = append(parts, slug)
	}
	body := strings.Join(parts, "-")

	t := normalizeType(taskType)
	if t != "" && body != "" {
		return t + "/" + body
	}
	if t != "" {
		return t
	}
	return body
}
