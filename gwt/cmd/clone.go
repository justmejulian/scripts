package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <url> [dir]",
	Short: "Clone a repo as bare for worktree usage",
	Long: `Clone a repository as a bare repo configured for worktree usage.

Creates a .bare directory with the git data, a .git file pointing
to it, and configures fetch for remote branch tracking.

Examples:
  gwt clone git@github.com:org/repo.git
  gwt clone https://github.com/org/repo.git my-repo
  cd $(gwt clone git@github.com:org/repo.git)`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

var repoNameRe = regexp.MustCompile(`([^/]+?)(?:\.git)?$`)

func extractRepoName(url string) (string, error) {
	// Handle SSH URLs like git@github.com:org/repo.git
	if i := strings.LastIndex(url, ":"); i != -1 && !strings.Contains(url, "://") {
		url = url[i+1:]
	}
	matches := repoNameRe.FindStringSubmatch(url)
	if matches == nil || matches[1] == "" {
		return "", fmt.Errorf("cannot extract repo name from %q", url)
	}
	return matches[1], nil
}

func runClone(cmd *cobra.Command, args []string) error {
	url := args[0]

	var dir string
	if len(args) > 1 {
		dir = args[1]
	} else {
		name, err := extractRepoName(url)
		if err != nil {
			return err
		}
		dir = name
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(absDir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	fmt.Fprintf(os.Stderr, "gwt clone: cloning %s...\n", url)
	if err := gitIn(absDir, "clone", "--bare", url, ".bare"); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(absDir, ".git"), []byte("gitdir: ./.bare\n"), 0o644); err != nil {
		return fmt.Errorf("writing .git file: %w", err)
	}

	fmt.Fprintln(os.Stderr, "gwt clone: configuring fetch...")
	if err := gitIn(absDir, "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*"); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "gwt clone: fetching...")
	if err := gitIn(absDir, "fetch", "origin"); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, absDir)
	return nil
}

func gitIn(dir string, args ...string) error {
	c := exec.Command("git", args...)
	c.Dir = dir
	c.Stderr = os.Stderr
	if out, err := c.Output(); err != nil {
		return fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(string(out)))
	}
	return nil
}
