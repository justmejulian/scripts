package cmd

import "testing"

func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"git@github.com:org/repo.git", "repo"},
		{"git@github.com:org/repo", "repo"},
		{"https://github.com/org/repo.git", "repo"},
		{"https://github.com/org/repo", "repo"},
		{"git@gitlab.com:group/subgroup/repo.git", "repo"},
		{"https://dev.azure.com/org/project/_git/repo", "repo"},
		{"/local/path/to/repo.git", "repo"},
		{"ssh://git@github.com/org/my-repo.git", "my-repo"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got, err := extractRepoName(tt.url)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("extractRepoName(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}
