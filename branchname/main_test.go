package main

import "testing"

func TestSlugifyDescription(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"add login page", "add-login-page"},
		{"Add Login Page", "add-login-page"},
		{"add - login page", "add-login-page"},
		{"add  multiple   spaces", "add-multiple-spaces"},
		{"hello, world!", "hello-world"},
		{"feat: something (cool)", "feat-something-cool"},
		{"already-slugged", "already-slugged"},
		{"--leading-trailing--", "leading-trailing"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugifyDescription(tt.input)
			if got != tt.want {
				t.Errorf("slugifyDescription(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildBranchName(t *testing.T) {
	tests := []struct {
		taskType    string
		issue       string
		description string
		want        string
	}{
		{"feat", "PROJ-123", "add login page", "feat/PROJ-123-add-login-page"},
		{"feat", "PROJ-123", "", "feat/PROJ-123"},
		{"feat", "", "add login page", "feat/add-login-page"},
		{"", "PROJ-123", "add login page", "PROJ-123-add-login-page"},
		{"", "PROJ-123", "", "PROJ-123"},
		{"", "", "add login page", "add-login-page"},
		{"", "", "", ""},
		{"FEAT", "PROJ-123", "add login page", "feat/PROJ-123-add-login-page"},
		{"chore", "PROJ-42", "update - dependencies", "chore/PROJ-42-update-dependencies"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := buildBranchName(tt.taskType, tt.issue, tt.description)
			if got != tt.want {
				t.Errorf("buildBranchName(%q, %q, %q) = %q, want %q",
					tt.taskType, tt.issue, tt.description, got, tt.want)
			}
		})
	}
}
