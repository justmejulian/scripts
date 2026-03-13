package main

import (
	"flag"
	"fmt"
	"os"

	branchutil "scripts/internal/branchname"
)

func slugifyDescription(desc string) string {
	return branchutil.SlugifyDescription(desc)
}

func buildBranchName(taskType, issue, description string) string {
	return branchutil.BuildName(taskType, issue, description)
}

func main() {
	taskType := flag.String("type", "", "task type (e.g. feat, fix, chore)")
	issue := flag.String("issue", "", "issue key (e.g. PROJ-123)")
	description := flag.String("description", "", "short description")

	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	result := buildBranchName(*taskType, *issue, *description)
	if result != "" {
		fmt.Println(result)
	}
}
