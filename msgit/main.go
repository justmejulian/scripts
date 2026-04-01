package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"scripts/internal/ai/providers/ollama"
	ollamamodels "scripts/internal/ai/providers/ollama/models"
	"scripts/internal/ai/providers/zen"
	zenmodels "scripts/internal/ai/providers/zen/models"
	"scripts/internal/ai/spec"
	"scripts/internal/ai/spec/model"

	"github.com/spf13/cobra"
)

var jiraRe = regexp.MustCompile(`[A-Z][A-Z0-9]+-\d+`)

type providerConfig struct {
	model       model.Info
	config      string
	newProvider func() (spec.Provider, error)
}

func selectedProviderConfig(offline bool) providerConfig {
	if offline {
		m := ollamamodels.Qwen3_5_4B
		return providerConfig{
			model:       m.Info,
			config:      m.Config.ThinkDisabled,
			newProvider: ollama.New,
		}
	}
	m := zenmodels.GPT5Nano
	return providerConfig{
		model:       m.Info,
		config:      m.Config.Default,
		newProvider: zen.New,
	}
}

var rootCmd = &cobra.Command{
	Use:           "msgit",
	Short:         "Generate a commit message for staged changes using AI",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          run,
}

func init() {
	rootCmd.Flags().Bool("offline", false, "use local ollama instead of zen")
}

func run(cmd *cobra.Command, args []string) error {
	offline, _ := cmd.Flags().GetBool("offline")

	providerCfg := selectedProviderConfig(offline)

	fmt.Fprintln(os.Stderr, "msgit: reading staged diff...")
	g := NewGit()
	diff, err := g.StagedDiff()
	if err != nil {
		return fmt.Errorf("msgit: failed to get staged diff: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("msgit: nothing staged (run git add first)")
	}

	branch := g.CurrentBranch()
	log, _ := g.RecentLog(5)

	prompt := buildPrompt(branch, strings.TrimSpace(log), diff)

	fmt.Fprintf(os.Stderr, "msgit: asking %s via %s...\n", providerCfg.model.Name, providerCfg.model.Provider)

	provider, err := providerCfg.newProvider()
	if err != nil {
		return fmt.Errorf("msgit: ai setup error: %w", err)
	}

	resp, err := provider.Generate(context.Background(), spec.Request{
		Prompt: prompt,
		Model:  providerCfg.model,
		Config: providerCfg.config,
	})
	fmt.Fprintln(os.Stderr, "")
	if err != nil {
		return fmt.Errorf("msgit: ai error: %w", err)
	}

	fmt.Print(resp.Text)
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func extractJiraKey(branch string) string {
	return jiraRe.FindString(strings.ToUpper(branch))
}

func buildPrompt(branch, log, diff string) string {
	jiraInstruction := ""
	if key := extractJiraKey(branch); key != "" {
		jiraInstruction = fmt.Sprintf("\n- Prepend the Jira issue key to the first line: %s <type>(<scope>): <summary>", key)
	}
	return fmt.Sprintf(`You are a commit message generator.

CRITICAL: Output ONLY the raw commit message text. Do NOT include:
- Any explanation or commentary
- Markdown formatting or code fences
- Headers, bullet points, or lists
- <think> tags or any XML tags
- Anything before or after the commit message

Commit message format:
- First line: imperative mood, max 72 chars, format: <type>(<scope>): <summary>
- Types: feat, fix, refactor, docs, test, chore
- Body (optional): separated by blank line, explain *why*%s

Context:
Branch: %s
Recent commits:
%s

Staged diff:
%s`, jiraInstruction, branch, log, diff)
}
