package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// CurrentBranch returns the name of the current git branch.
func CurrentBranch() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("git: not in a git repository or git not found: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
