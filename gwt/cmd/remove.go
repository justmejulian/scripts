package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"scripts/internal/git"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <path-or-name>",
	Short: "Remove a worktree",
	Long: `Remove a worktree, prune stale references, and delete the branch.

Accepts an absolute path or a directory name (resolved against
the bare repo root). Uses safe branch deletion (-d); unmerged
branches are kept with a warning.

Examples:
  gwt remove my-feature
  gwt remove /absolute/path/to/worktree
  gwt remove $(gwt list | fzf)
  gwt remove -f dirty-worktree`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	RunE:    runRemove,
}

var removeForce bool

func init() {
	removeCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "force remove dirty worktrees")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	target := args[0]

	// Resolve to absolute path
	if !filepath.IsAbs(target) {
		root, err := git.Root()
		if err != nil {
			return err
		}
		target = filepath.Join(root, target)
	}

	// Safety: refuse to remove cwd
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}
	if strings.HasPrefix(cwd, target) {
		return fmt.Errorf("cannot remove the current worktree; switch to another first")
	}

	// Find the branch for this worktree
	branch := ""
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return err
	}
	for _, wt := range worktrees {
		if wt.Path == target {
			branch = wt.Branch
			break
		}
	}

	fmt.Fprintf(os.Stderr, "gwt remove: removing worktree at %s...\n", target)
	if err := git.RemoveWorktree(target, removeForce); err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "gwt remove: pruning stale worktrees...")
	if err := git.PruneWorktrees(); err != nil {
		return err
	}

	if branch != "" {
		fmt.Fprintf(os.Stderr, "gwt remove: deleting branch %s...\n", branch)
		if err := git.DeleteBranch(branch); err != nil {
			fmt.Fprintf(os.Stderr, "gwt remove: warning: %v\n", err)
		}
	}

	fmt.Fprintln(os.Stdout, target)
	return nil
}
