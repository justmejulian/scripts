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

func main() {
	fmt.Fprintln(os.Stderr, "msgit: reading staged diff...")
	diff, err := stagedDiff()
	if err != nil {
		fmt.Fprintln(os.Stderr, "msgit: failed to get staged diff:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(diff) == "" {
		fmt.Fprintln(os.Stderr, "msgit: nothing staged (run git add first)")
		os.Exit(1)
	}

	branch := currentBranch()
	log, _ := recentLog(5)

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

func buildPrompt(branch, log, diff string) string {
	return fmt.Sprintf(`You are a commit message generator. Output ONLY the commit message — no explanation, no markdown fences, no <think> tags.

Conventions:
- First line: imperative mood, max 72 chars, format: <type>(<scope>): <summary>
- Types: feat, fix, refactor, docs, test, chore
- Body (optional): separated by blank line, explain *why*

Context:
Branch: %s
Recent commits:
%s

Staged diff:
%s`, branch, log, diff)
}

func stripThinking(s string) string {
	return strings.TrimSpace(thinkRe.ReplaceAllString(s, ""))
}
