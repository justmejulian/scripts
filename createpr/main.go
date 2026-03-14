package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"scripts/internal/azure"
	"scripts/internal/branchname"
	"scripts/internal/git"
	"scripts/internal/jira"
)

var jiraKeyRe = regexp.MustCompile(`[A-Z]+-[0-9]+`)

func main() {
	target := flag.String("target", "main", "target branch for the PR")

	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	ctx := context.Background()

	// Derive project + repo from working directory
	wd, err := os.Getwd()
	if err != nil {
		fail(fmt.Errorf("could not get working directory: %w", err))
	}
	parts := strings.Split(filepath.ToSlash(wd), "/")
	if len(parts) < 2 {
		fail(fmt.Errorf("working directory %q has fewer than 2 path segments", wd))
	}
	project := parts[len(parts)-2]
	repo := parts[len(parts)-1]

	// Get current branch
	branch, err := git.CurrentBranch()
	if err != nil {
		fail(err)
	}

	// Extract branch type
	branchType, err := branchname.BranchType(branch)
	if err != nil {
		fail(fmt.Errorf("could not determine branch type: branch must contain '/' (e.g. feat/PROJ-123-...)"))
	}

	// Extract Jira issue key
	key := jiraKeyRe.FindString(branch)
	if key == "" {
		fail(fmt.Errorf("no Jira issue key found in branch %q (expected pattern like PROJ-123)", branch))
	}

	// Build clients
	jiraClient, err := jira.NewClientFromEnv()
	if err != nil {
		fail(err)
	}
	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		fail(err)
	}

	// Get issue title
	issue, err := jiraClient.GetIssue(ctx, key)
	if err != nil {
		fail(fmt.Errorf("could not fetch Jira issue %s: %w", key, err))
	}

	// Build PR title
	prTitle := fmt.Sprintf("%s %s: %s", key, branchType, issue.Title)

	// Create PR
	pr, err := azureClient.CreatePR(ctx, project, repo, azure.CreatePRRequest{
		Title:         prTitle,
		SourceRefName: "refs/heads/" + branch,
		TargetRefName: "refs/heads/" + *target,
	})
	if err != nil {
		fail(fmt.Errorf("could not create PR: %w", err))
	}

	org := os.Getenv("AZURE_DEVOPS_ORG")
	prURL := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d", org, project, repo, pr.PullRequestID)
	fmt.Println(prURL)

	// Transition Jira issue to "In Review"
	if err := jiraClient.TransitionIssue(ctx, key, "In Review"); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not transition %s to 'In Review': %v\n", key, err)
	}

	// Add PR URL as Jira comment
	if err := jiraClient.AddComment(ctx, key, prURL); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not add comment to %s: %v\n", key, err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
