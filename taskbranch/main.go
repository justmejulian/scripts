package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	branchutil "scripts/internal/branchname"
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

	reader := bufio.NewReader(os.Stdin)
	issue, err := selectIssue(reader, issues)
	if err != nil {
		fail(err)
	}

	selectedType := strings.TrimSpace(*taskType)
	if selectedType == "" {
		selectedType, err = selectBranchType(reader)
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

func selectIssue(reader *bufio.Reader, issues []jira.Issue) (jira.Issue, error) {
	for {
		fmt.Fprintln(os.Stderr, "Assigned tasks:")
		for i, issue := range issues {
			label := fmt.Sprintf("%d. %s - %s", i+1, issue.Key, issue.Title)
			if strings.TrimSpace(issue.Status) != "" {
				label += " [" + issue.Status + "]"
			}
			fmt.Fprintln(os.Stderr, label)
		}

		fmt.Fprintf(os.Stderr, "Select task [1-%d]: ", len(issues))
		input, err := readLine(reader)
		if err != nil {
			return jira.Issue{}, err
		}

		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(issues) {
			fmt.Fprintln(os.Stderr, "Invalid selection, try again.")
			continue
		}

		return issues[index-1], nil
	}
}

func selectBranchType(reader *bufio.Reader) (string, error) {
	for {
		fmt.Fprintln(os.Stderr, "Branch types:")
		for i, branchType := range branchTypes {
			fmt.Fprintf(os.Stderr, "%d. %s\n", i+1, branchType)
		}

		fmt.Fprintf(os.Stderr, "Select branch type [1-%d]: ", len(branchTypes))
		input, err := readLine(reader)
		if err != nil {
			return "", err
		}

		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(branchTypes) {
			fmt.Fprintln(os.Stderr, "Invalid selection, try again.")
			continue
		}

		selected := branchTypes[index-1]
		if selected != "custom" {
			return selected, nil
		}

		fmt.Fprint(os.Stderr, "Enter custom branch type: ")
		customType, err := readLine(reader)
		if err != nil {
			return "", err
		}

		customType = strings.TrimSpace(customType)
		if customType == "" {
			fmt.Fprintln(os.Stderr, "Branch type cannot be empty.")
			continue
		}

		return customType, nil
	}
}

func readLine(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			return strings.TrimSpace(input), nil
		}
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
