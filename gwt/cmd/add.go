package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"scripts/internal/git"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <branch>",
	Short: "Add a worktree for a branch",
	Long: `Add a worktree for a branch and print its path.

The worktree directory is created under the bare repo root, using
the last component of the branch name as the directory name
(e.g. feat/PROJ-123-foo becomes PROJ-123-foo).

Examples:
  gwt add main
  gwt add -b feat/PROJ-123-new-feature
  cd $(gwt add -b $(taskbranch))`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

var addNewBranch bool

func init() {
	addCmd.Flags().BoolVarP(&addNewBranch, "branch", "b", false, "create a new branch")
	rootCmd.AddCommand(addCmd)
}

// worktreeDirName returns the directory name for a worktree.
// It strips any prefix before the last "/" to avoid nested directories.
func worktreeDirName(branch string) string {
	if i := strings.LastIndex(branch, "/"); i != -1 {
		return branch[i+1:]
	}
	return branch
}

func runAdd(cmd *cobra.Command, args []string) error {
	branch := args[0]

	root, err := git.Root()
	if err != nil {
		return err
	}

	dir := worktreeDirName(branch)
	path := filepath.Join(root, dir)

	if addNewBranch {
		err = git.CreateWorktree(path, branch)
	} else {
		err = git.AddWorktree(path, branch)
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "gwt add: created worktree for %s\n", branch)
	fmt.Fprintln(os.Stdout, path)
	return nil
}
