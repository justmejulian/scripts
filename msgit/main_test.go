package main

import (
	"strings"
	"testing"
)

func TestStripThinking(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no think tags",
			input: "feat(scope): summary",
			want:  "feat(scope): summary",
		},
		{
			name:  "think tag stripped",
			input: "<think>internal reasoning here</think>\nfeat(scope): summary",
			want:  "feat(scope): summary",
		},
		{
			name:  "multiline think tag stripped",
			input: "<think>\nline one\nline two\n</think>\nfix(auth): correct token expiry",
			want:  "fix(auth): correct token expiry",
		},
		{
			name:  "leading and trailing whitespace trimmed",
			input: "  \n<think>x</think>\n  chore: update deps  \n",
			want:  "chore: update deps",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only think tag",
			input: "<think>nothing useful</think>",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripThinking(tt.input)
			if got != tt.want {
				t.Errorf("stripThinking(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExtractJiraKey(t *testing.T) {
	tests := []struct {
		branch string
		want   string
	}{
		{"feat/SA-2809-cell-status-card", "SA-2809"},
		{"feat/PROJ-123-some-feature", "PROJ-123"},
		{"main", ""},
		{"feature/no-key-here", ""},
		{"fix/ABC-1-short", "ABC-1"},
	}
	for _, tt := range tests {
		got := extractJiraKey(tt.branch)
		if got != tt.want {
			t.Errorf("extractJiraKey(%q) = %q, want %q", tt.branch, got, tt.want)
		}
	}
}

func TestBuildPrompt(t *testing.T) {
	branch := "feature/login"
	log := "abc1234 feat(auth): add JWT support"
	diff := "diff --git a/main.go b/main.go\n+added line"

	prompt := buildPrompt(branch, log, diff)

	for _, want := range []string{branch, log, diff} {
		if !strings.Contains(prompt, want) {
			t.Errorf("buildPrompt output missing %q", want)
		}
	}

	if !strings.Contains(prompt, "commit message") {
		t.Error("buildPrompt output missing instructions")
	}
}

func TestBuildPromptWithJiraKey(t *testing.T) {
	prompt := buildPrompt("feat/SA-2809-cell-status-card", "", "diff --git a/x b/x\n+line")
	if !strings.Contains(prompt, "SA-2809") {
		t.Error("buildPrompt missing Jira key instruction")
	}
}

func TestBuildPromptWithoutJiraKey(t *testing.T) {
	prompt := buildPrompt("feature/no-ticket", "", "diff --git a/x b/x\n+line")
	if strings.Contains(prompt, "Jira") {
		t.Error("buildPrompt should not mention Jira when no key in branch")
	}
}
