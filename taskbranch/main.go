package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	branchutil "scripts/internal/branchname"
	"scripts/internal/fzf"
	"scripts/internal/jira"

	"github.com/spf13/cobra"
)

const defaultJQL = "assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC"

var branchTypes = []string{"feat", "fix", "chore", "custom"}

var rootCmd = &cobra.Command{
	Use:           "taskbranch",
	Short:         "Interactively select a Jira task and generate a branch name",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func init() {
	rootCmd.Flags().String("type", "", "branch type (e.g. feat, fix, chore); skips prompt when set")
	rootCmd.Flags().String("jql", defaultJQL, "Jira query used to fetch assigned tasks")
}

func run(cmd *cobra.Command, args []string) error {
	taskType, _ := cmd.Flags().GetString("type")
	jql, _ := cmd.Flags().GetString("jql")

	client, err := jira.NewClientFromEnv()
	if err != nil {
		return err
	}

	issues, err := client.SearchIssues(context.Background(), jql, []string{"summary", "status"})
	if err != nil {
		return err
	}

	if len(issues) == 0 {
		return fmt.Errorf("no assigned tasks found for query: %s", jql)
	}

	issue, err := selectIssue(issues)
	if err != nil {
		return err
	}

	selectedType := strings.TrimSpace(taskType)
	if selectedType == "" {
		selectedType, err = selectBranchType()
		if err != nil {
			return err
		}
	}

	result := branchutil.BuildName(selectedType, issue.Key, issue.Title)
	if result == "" {
		return fmt.Errorf("could not build branch name from selected task")
	}

	fmt.Println(result)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
