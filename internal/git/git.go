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

// RepoRoot returns the absolute path of the top-level git repository directory.
func RepoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git: not in a git repository or git not found: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// PushBranch pushes the given branch to origin, setting the upstream.
func PushBranch(branch string) error {
	cmd := exec.Command("git", "push", "-u", "origin", branch)
	cmd.Stdout = nil
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git push failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
