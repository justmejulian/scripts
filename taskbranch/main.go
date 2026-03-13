package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	branchutil "scripts/internal/branchname"
	"scripts/internal/fzf"
	"scripts/internal/jira"
)

const defaultJQL = "assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC"

var branchTypes = []string{"feat", "fix", "chore", "custom"}

func main() {
	taskType := flag.String("type", "", "branch type (e.g. feat, fix, chore); skips prompt when set")
	jql := flag.String("jql", defaultJQL, "Jira query used to fetch assigned tasks")

	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	client, err := jira.NewClientFromEnv()
	if err != nil {
		fail(err)
	}

	issues, err := client.SearchIssues(context.Background(), *jql, []string{"summary", "status"})
	if err != nil {
		fail(err)
	}

	if len(issues) == 0 {
		fail(fmt.Errorf("no assigned tasks found for query: %s", *jql))
	}

	issue, err := selectIssue(issues)
	if err != nil {
		fail(err)
	}

	selectedType := strings.TrimSpace(*taskType)
	if selectedType == "" {
		selectedType, err = selectBranchType()
		if err != nil {
			fail(err)
		}
	}

	result := branchutil.BuildName(selectedType, issue.Key, issue.Title)
	if result == "" {
		fail(fmt.Errorf("could not build branch name from selected task"))
	}

	fmt.Println(result)
}

func selectIssue(issues []jira.Issue) (jira.Issue, error) {
	lines := make([]string, len(issues))
	for i, issue := range issues {
		label := fmt.Sprintf("%s - %s", issue.Key, issue.Title)
		if strings.TrimSpace(issue.Status) != "" {
			label += " [" + issue.Status + "]"
		}
		lines[i] = label
	}

	result, err := fzf.Select("Task:", lines)
	if err != nil {
		if errors.Is(err, fzf.ErrCancelled) {
			return jira.Issue{}, fmt.Errorf("no task selected")
		}
		return jira.Issue{}, err
	}

	key := strings.SplitN(result, " - ", 2)[0]
	for _, issue := range issues {
		if issue.Key == key {
			return issue, nil
		}
	}

	return jira.Issue{}, fmt.Errorf("selected issue not found: %s", key)
}

func selectBranchType() (string, error) {
	result, err := fzf.Select("Branch type:", branchTypes)
	if err != nil {
		if errors.Is(err, fzf.ErrCancelled) {
			return "", fmt.Errorf("no branch type selected")
		}
		return "", err
	}

	if result != "custom" {
		return result, nil
	}

	fmt.Fprint(os.Stderr, "Enter custom branch type: ")
	reader := bufio.NewReader(os.Stdin)
	customType, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	customType = strings.TrimSpace(customType)
	if customType == "" {
		return "", fmt.Errorf("branch type cannot be empty")
	}

	return customType, nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
