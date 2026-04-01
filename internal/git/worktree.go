package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a single git worktree entry.
type Worktree struct {
	Path   string
	Branch string
	IsBare bool
}

// ListWorktrees parses the output of git worktree list --porcelain.
func ListWorktrees() ([]Worktree, error) {
	out, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}
	return parseWorktreeList(string(out)), nil
}

func parseWorktreeList(output string) []Worktree {
	var worktrees []Worktree
	var current *Worktree

	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			if current != nil {
				worktrees = append(worktrees, *current)
			}
			current = &Worktree{Path: strings.TrimPrefix(line, "worktree ")}
		} else if line == "bare" && current != nil {
			current.IsBare = true
		} else if strings.HasPrefix(line, "branch ") && current != nil {
			ref := strings.TrimPrefix(line, "branch ")
			current.Branch = strings.TrimPrefix(ref, "refs/heads/")
		}
	}
	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees
}

// Root returns the root directory of a bare worktree setup.
// It uses git rev-parse --git-common-dir to find the .bare directory,
// then returns its parent.
func Root() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--git-common-dir").Output()
	if err != nil {
		return "", fmt.Errorf("git: not in a git repository: %w", err)
	}
	commonDir := strings.TrimSpace(string(out))

	abs, err := filepath.Abs(commonDir)
	if err != nil {
		return "", fmt.Errorf("git: resolving path: %w", err)
	}

	return filepath.Dir(abs), nil
}

// AddWorktree checks out an existing branch into a new worktree.
func AddWorktree(path, branch string) error {
	out, err := exec.Command("git", "worktree", "add", path, branch).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// CreateWorktree creates a new branch and checks it out into a new worktree.
func CreateWorktree(path, branch string) error {
	out, err := exec.Command("git", "worktree", "add", "-b", branch, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add -b: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// RemoveWorktree removes a worktree directory.
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove", path}
	if force {
		args = append(args, "--force")
	}
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree remove: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// PruneWorktrees cleans up stale worktree references.
func PruneWorktrees() error {
	out, err := exec.Command("git", "worktree", "prune").CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree prune: %s", strings.TrimSpace(string(out)))
	}
	return nil
}

// DeleteBranch deletes a local branch. Uses -d (safe delete).
func DeleteBranch(branch string) error {
	out, err := exec.Command("git", "branch", "-d", branch).CombinedOutput()
	if err != nil {
		return fmt.Errorf("git branch -d: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
