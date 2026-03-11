package main

import (
	"context"
	"fmt"
	"os"

	"scripts/internal/jira"
)

func main() {
	c, err := jira.NewClientFromEnv()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx := context.Background()
	issues, err := c.SearchIssues(ctx, "assignee = currentUser() AND status = 3", []string{"key", "summary", "status"})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, issue := range issues {
		fmt.Printf("%s: %s\n", issue.Key, issue.Title)
	}
}
