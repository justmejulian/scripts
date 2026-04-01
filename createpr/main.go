package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"scripts/internal/azure"
	"scripts/internal/branchname"
	"scripts/internal/git"
	"scripts/internal/jira"
	"scripts/internal/prompt"

	"github.com/spf13/cobra"
)

var jiraKeyRe = regexp.MustCompile(`[A-Z]+-[0-9]+`)

var rootCmd = &cobra.Command{
	Use:           "createpr",
	Short:         "Create a pull request in Azure DevOps for the current branch",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func init() {
	rootCmd.Flags().String("target", "main", "target branch for the PR")
}

func run(cmd *cobra.Command, args []string) error {
	target, _ := cmd.Flags().GetString("target")

	ctx := context.Background()

	project, repo, err := resolveRepoContext()
	if err != nil {
		return err
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	branchType, key, err := parseBranch(branch)
	if err != nil {
		return err
	}

	jiraClient, err := jira.NewClientFromEnv()
	if err != nil {
		return err
	}
	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		return err
	}

	issue, err := jiraClient.GetIssue(ctx, key)
	if err != nil {
		return fmt.Errorf("could not fetch Jira issue %s: %w", key, err)
	}

	req := azure.CreatePRRequest{
		Title:         buildPRTitle(key, branchType, issue.Title),
		SourceRefName: "refs/heads/" + branch,
		TargetRefName: "refs/heads/" + target,
	}
	pr, err := azureClient.CreatePR(ctx, project, repo, req)
	if err != nil {
		var apiErr *azure.APIError
		if errors.As(err, &apiErr) && strings.Contains(apiErr.Body, "TF401398") {
			if prompt.Confirm("branch not found on remote — push it now? [y/N]: ") {
				if pushErr := git.PushBranch(branch); pushErr != nil {
					return fmt.Errorf("could not push branch: %w", pushErr)
				}
				pr, err = azureClient.CreatePR(ctx, project, repo, req)
			}
		}
		if err != nil {
			return fmt.Errorf("could not create PR: %w", err)
		}
	}

	prURL := buildPRURL(os.Getenv("AZURE_DEVOPS_ORG"), project, repo, pr.PullRequestID)
	fmt.Println(prURL)

	updateJiraAfterPR(ctx, jiraClient, key, prURL)

	jiraURL := jiraClient.BrowseURL(key)
	fmt.Printf("\nPR for <%s|%s> is ready for review\n%s\n", jiraURL, key, prURL)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func resolveRepoContext() (project, repo string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("could not get working directory: %w", err)
	}
	parts := strings.Split(filepath.ToSlash(wd), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("working directory %q has fewer than 2 path segments", wd)
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}

func parseBranch(branch string) (branchType, jiraKey string, err error) {
	branchType, err = branchname.BranchType(branch)
	if err != nil {
		return "", "", fmt.Errorf("could not determine branch type: branch must contain '/' (e.g. feat/PROJ-123-...)")
	}
	jiraKey = jiraKeyRe.FindString(branch)
	if jiraKey == "" {
		return "", "", fmt.Errorf("no Jira issue key found in branch %q (expected pattern like PROJ-123)", branch)
	}
	return branchType, jiraKey, nil
}

func buildPRTitle(jiraKey, branchType, issueTitle string) string {
	return fmt.Sprintf("%s %s: %s", jiraKey, branchType, issueTitle)
}

func buildPRURL(org, project, repo string, prID int) string {
	return fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d", org, project, repo, prID)
}

func updateJiraAfterPR(ctx context.Context, client *jira.Client, key, prURL string) {
	if err := client.TransitionIssue(ctx, key, "In Review"); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not transition %s to 'In Review': %v\n", key, err)
	}
	if err := client.AddComment(ctx, key, "PR: "+prURL); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not add comment to %s: %v\n", key, err)
	}
	if err := client.UpdateIssue(ctx, key, map[string]any{"assignee": nil}); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not unassign %s: %v\n", key, err)
	}
}
