package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"scripts/internal/azure"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "revu",
	Short:         "CLI PR review tool",
	SilenceUsage:  true,
	SilenceErrors: true,
}

var commentsCmd = &cobra.Command{
	Use:   "comments <pr-id>",
	Short: "Display comments for a PR",
	Args:  cobra.ExactArgs(1),
	RunE:  runComments,
}

func init() {
	commentsCmd.Flags().String("project", "", "Azure DevOps project (required)")
	commentsCmd.Flags().String("repo", "", "Azure DevOps repository (required)")
	commentsCmd.MarkFlagRequired("project")
	commentsCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(commentsCmd)
}

func runComments(cmd *cobra.Command, args []string) error {
	prID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("pr-id must be an integer: %w", err)
	}

	project, _ := cmd.Flags().GetString("project")
	repo, _ := cmd.Flags().GetString("repo")

	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		return err
	}

	provider := &azureProvider{client: azureClient, project: project, repo: repo}

	threads, err := provider.GetThreads(context.Background(), prID)
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
