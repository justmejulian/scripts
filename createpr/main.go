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
	"scripts/internal/slack"
)

var jiraKeyRe = regexp.MustCompile(`[A-Z]+-[0-9]+`)

func main() {
	target := flag.String("target", "main", "target branch for the PR")
	slackChannel := flag.String("slack-channel", "", "slack channel to notify when PR is ready (required)")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	if *slackChannel == "" {
		fmt.Fprintln(os.Stderr, "flag --slack-channel is required")
		os.Exit(1)
	}

	ctx := context.Background()

	project, repo, err := resolveRepoContext()
	if err != nil {
		fail(err)
	}

	branch, err := git.CurrentBranch()
	if err != nil {
		fail(err)
	}

	branchType, key, err := parseBranch(branch)
	if err != nil {
		fail(err)
	}

	jiraClient, err := jira.NewClientFromEnv()
	if err != nil {
		fail(err)
	}
	azureClient, err := azure.NewClientFromEnv()
	if err != nil {
		fail(err)
	}
	slackClient, err := slack.NewClientFromEnv()
	if err != nil {
		fail(err)
	}

	issue, err := jiraClient.GetIssue(ctx, key)
	if err != nil {
		fail(fmt.Errorf("could not fetch Jira issue %s: %w", key, err))
	}

	pr, err := azureClient.CreatePR(ctx, project, repo, azure.CreatePRRequest{
		Title:         buildPRTitle(key, branchType, issue.Title),
		SourceRefName: "refs/heads/" + branch,
		TargetRefName: "refs/heads/" + *target,
	})
	if err != nil {
		fail(fmt.Errorf("could not create PR: %w", err))
	}

	prURL := buildPRURL(os.Getenv("AZURE_DEVOPS_ORG"), project, repo, pr.PullRequestID)
	fmt.Println(prURL)

	updateJiraAfterPR(ctx, jiraClient, key, prURL)

	msg := fmt.Sprintf("PR for %s is ready for review\n%s", key, prURL)
	if err := slackClient.PostMessage(ctx, *slackChannel, msg); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not send Slack message: %v\n", err)
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
	if err := client.AddComment(ctx, key, prURL); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not add comment to %s: %v\n", key, err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
