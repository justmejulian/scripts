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

var syncClean bool
var syncActiveOnly bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Inject PR comments into source files as code comments",
	Args:  cobra.NoArgs,
	RunE:  runSync,
}

func init() {
	rootCmd.AddCommand(commentsCmd)
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVar(&syncClean, "clean", false, "remove injected REVU comments without re-inserting")
	syncCmd.Flags().BoolVar(&syncActiveOnly, "active-only", false, "only sync active (unresolved) threads")
}

func fetchThreads(ctx context.Context) ([]Thread, error) {
	project, repo, err := repocontext.Resolve()
	if err != nil {
		return nil, err
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		return nil, err
	}

	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	pr, err := azureClient.GetPRByBranch(ctx, project, repo, branch)
	if err != nil {
		return nil, fmt.Errorf("could not find PR: %w", err)
	}

	provider := &azureProvider{client: azureClient, project: project, repo: repo}

	threads, err := provider.GetThreads(ctx, pr.PullRequestID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch threads: %w", err)
	}

	return threads, nil
}

func runComments(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	threads, err := fetchThreads(ctx)
	if err != nil {
		return err
	}

	if len(threads) == 0 {
		fmt.Println("no comments")
		return nil
	}

	printThreads(os.Stdout, threads)
	return nil
}

func runSync(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	threads, err := fetchThreads(ctx)
	if err != nil {
		return err
	}

	repoRoot, err := git.RepoRoot()
	if err != nil {
		return err
	}

	return syncFiles(threads, repoRoot, syncClean, syncActiveOnly)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
