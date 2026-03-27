package main

import (
	"strings"
	"testing"
)

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
