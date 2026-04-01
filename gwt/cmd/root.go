package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gwt",
	Short: "Git worktree manager for bare repos",
	Long: `gwt manages git worktrees in a bare repository setup.

It expects a repo cloned with "gwt clone", which creates a .bare
directory and a .git file pointing to it. Worktrees are added as
sibling directories next to .bare.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
