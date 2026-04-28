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

var syncActiveOnly bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Inject PR comments into source files as code comments",
	Args:  cobra.NoArgs,
	RunE:  runSync,
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all injected REVU comments (including REVU[NEW]) from source files",
	Args:  cobra.NoArgs,
	RunE:  runClean,
}

var uploadContextLines int
var uploadDryRun bool

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Interactively review and upload REVU[NEW] comments as PR threads",
	Args:  cobra.NoArgs,
	RunE:  runUpload,
}

func init() {
	rootCmd.AddCommand(commentsCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(uploadCmd)
	syncCmd.Flags().BoolVar(&syncActiveOnly, "active-only", false, "only sync active (unresolved) threads")
	uploadCmd.Flags().IntVar(&uploadContextLines, "context", 4, "lines of context to show around each comment")
	uploadCmd.Flags().BoolVar(&uploadDryRun, "dry-run", false, "show pending comments without uploading")
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

func runUpload(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repoRoot, err := git.RepoRoot()
	if err != nil {
		return err
	}

	pending, err := scanNewComments(repoRoot)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		fmt.Println("no REVU[NEW] comments found")
		return nil
	}

	fileCount := make(map[string]struct{})
	for _, c := range pending {
		fileCount[c.AbsPath] = struct{}{}
	}
	fmt.Printf("Found %d REVU[NEW] comment(s) across %d file(s).\n", len(pending), len(fileCount))

	if uploadDryRun {
		for _, c := range pending {
			fmt.Printf("\n--- %s:%d ---\n\n", c.RepoPath, c.Line)
			fmt.Print(renderContext(c.AbsPath, c.Line, uploadContextLines))
		}
		return nil
	}

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

	pr, err := azureClient.GetPRByBranch(ctx, project, repo, branch)
	if err != nil {
		return fmt.Errorf("could not find PR: %w", err)
	}

	provider := &azureProvider{client: azureClient, project: project, repo: repo}

	approved, err := interactiveReview(pending, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	if len(approved) == 0 {
		fmt.Println("nothing to upload")
		return nil
	}

	applyUploads(ctx, provider, pr.PullRequestID, approved, os.Stdout)
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

	return syncFiles(threads, repoRoot, false, syncActiveOnly)
}

func runClean(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	threads, err := fetchThreads(ctx)
	if err != nil {
		return err
	}

	repoRoot, err := git.RepoRoot()
	if err != nil {
		return err
	}

	return syncFiles(threads, repoRoot, true, false)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
