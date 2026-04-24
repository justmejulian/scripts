// Package repocontext derives the Azure DevOps project and repository names
// from the current working directory path.
package repocontext

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Resolve returns the Azure DevOps project and repository names derived
// from the last two segments of the current working directory path.
func Resolve() (project, repo string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("could not get working directory: %w", err)
	}
	parts := strings.Split(filepath.ToSlash(wd), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("working directory %q has fewer than 2 path segments", wd)
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}
