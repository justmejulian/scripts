package main

import (
	"context"
	"fmt"
	"os"

	"scripts/internal/azure"
	"scripts/internal/git"
	"scripts/internal/repocontext"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "revu",
	Short:         "CLI PR review tool",
	SilenceUsage:  true,
	SilenceErrors: true,
}

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Display comments for the PR on the current branch",
	Args:  cobra.NoArgs,
	RunE:  runComments,
}

func init() {
	rootCmd.AddCommand(commentsCmd)
}

func runComments(cmd *cobra.Command, args []string) error {
	project, repo, err := repocontext.Resolve()
	if err != nil {
		return err
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		return err
	}

	ctx := context.Background()

	pr, err := azureClient.GetPRByBranch(ctx, project, repo, branch)
	if err != nil {
		return fmt.Errorf("could not find PR: %w", err)
	}

	provider := &azureProvider{client: azureClient, project: project, repo: repo}

	threads, err := provider.GetThreads(ctx, pr.PullRequestID)
	if err != nil {
		return fmt.Errorf("could not fetch threads: %w", err)
	}

	if len(threads) == 0 {
		fmt.Println("no comments")
		return nil
	}

	printThreads(os.Stdout, threads)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
