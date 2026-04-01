package cmd

import (
	"os/exec"
	"strings"

	"scripts/internal/git"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate a shell completion script for gwt.

To load completions:

  bash:
    source <(gwt completion bash)

  zsh:
    source <(gwt completion zsh)
    # Or add to fpath for persistent completions:
    gwt completion zsh > "${fpath[1]}/_gwt"

  fish:
    gwt completion fish | source
    # Or persist:
    gwt completion fish > ~/.config/fish/completions/gwt.fish`,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Custom completions for add: complete branch names
	addCmd.ValidArgsFunction = completeBranches

	// Custom completions for remove: complete worktree names
	removeCmd.ValidArgsFunction = completeWorktrees
}

func completeBranches(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	out, err := exec.Command("git", "branch", "-a", "--format=%(refname:short)").Output()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var branches []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			branches = append(branches, line)
		}
	}
	return branches, cobra.ShellCompDirectiveNoFileComp
}

func completeWorktrees(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	worktrees, err := git.ListWorktrees()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, wt := range worktrees {
		if wt.IsBare {
			continue
		}
		names = append(names, wt.Path)
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}
