package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"scripts/internal/ollama"
)

var thinkRe = regexp.MustCompile(`(?s)<think>.*?</think>`)
var jiraRe = regexp.MustCompile(`[A-Z][A-Z0-9]+-\d+`)

func main() {
	fmt.Fprintln(os.Stderr, "msgit: reading staged diff...")
	g := NewGit()
	diff, err := g.StagedDiff()
	if err != nil {
		fmt.Fprintln(os.Stderr, "msgit: failed to get staged diff:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(diff) == "" {
		fmt.Fprintln(os.Stderr, "msgit: nothing staged (run git add first)")
		os.Exit(1)
	}

	branch := g.CurrentBranch()
	log, _ := g.RecentLog(5)

	prompt := buildPrompt(branch, strings.TrimSpace(log), diff)

	fmt.Fprintf(os.Stderr, "msgit: asking %s...\n", ollama.ModelQwen3_8B)

	c := ollama.NewClient(ollama.ModelQwen3_8B)
	reply, err := c.QuickChat(context.Background(), prompt)
	fmt.Fprintln(os.Stderr, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "msgit: ollama error:", err)
		os.Exit(1)
	}

	fmt.Print(stripThinking(reply))
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

func stripThinking(s string) string {
	return strings.TrimSpace(thinkRe.ReplaceAllString(s, ""))
}
