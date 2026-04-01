package cmd

import (
	"fmt"
	"os"

	"scripts/internal/git"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List worktree paths",
	Long: `List worktree paths, one per line.

Outputs non-bare worktree paths to stdout, suitable for piping
to fzf, grep, wc, or other tools.

Examples:
  gwt list                    # one path per line
  gwt list -b                 # path<tab>branch per line
  gwt list | fzf              # interactive selection
  gwt list -b | column -t     # pretty print`,
	Args: cobra.NoArgs,
	RunE: runList,
}

var listBranch bool

func init() {
	listCmd.Flags().BoolVarP(&listBranch, "branch", "b", false, "include branch name (tab-separated)")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return err
	}

	for _, wt := range worktrees {
		if wt.IsBare {
			continue
		}
		if listBranch {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", wt.Path, wt.Branch)
		} else {
			fmt.Fprintln(os.Stdout, wt.Path)
		}
	}
	return nil
}
