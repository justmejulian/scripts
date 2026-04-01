package main

import (
	"fmt"
	"os"

	branchutil "scripts/internal/branchname"

	"github.com/spf13/cobra"
)

func slugifyDescription(desc string) string {
	return branchutil.SlugifyDescription(desc)
}

func buildBranchName(taskType, issue, description string) string {
	return branchutil.BuildName(taskType, issue, description)
}

var rootCmd = &cobra.Command{
	Use:           "branchname",
	Short:         "Generate a git branch name from components",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func init() {
	rootCmd.Flags().String("type", "", "task type (e.g. feat, fix, chore)")
	rootCmd.Flags().String("issue", "", "issue key (e.g. PROJ-123)")
	rootCmd.Flags().String("description", "", "short description")
}

func run(cmd *cobra.Command, args []string) error {
	taskType, _ := cmd.Flags().GetString("type")
	issue, _ := cmd.Flags().GetString("issue")
	description, _ := cmd.Flags().GetString("description")

	result := buildBranchName(taskType, issue, description)
	if result != "" {
		fmt.Println(result)
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
