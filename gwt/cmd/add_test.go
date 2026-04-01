package cmd

import "testing"

func TestWorktreeDirName(t *testing.T) {
	tests := []struct {
		branch string
		want   string
	}{
		{"main", "main"},
		{"feat/PROJ-123-foo", "PROJ-123-foo"},
		{"fix/bug", "bug"},
		{"release/v1.0.0", "v1.0.0"},
		{"feature/nested/deep/branch", "branch"},
	}

	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			got := worktreeDirName(tt.branch)
			if got != tt.want {
				t.Errorf("worktreeDirName(%q) = %q, want %q", tt.branch, got, tt.want)
			}
		})
	}
}
