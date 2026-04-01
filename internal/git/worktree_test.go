package git

import "testing"

func TestParseWorktreeList(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []Worktree
	}{
		{
			name: "bare repo with two worktrees",
			input: `worktree /repo/.bare
HEAD abc1234
bare

worktree /repo/main
HEAD def5678
branch refs/heads/main

worktree /repo/feature
HEAD 9876543
branch refs/heads/feat/my-feature
`,
			want: []Worktree{
				{Path: "/repo/.bare", IsBare: true},
				{Path: "/repo/main", Branch: "main"},
				{Path: "/repo/feature", Branch: "feat/my-feature"},
			},
		},
		{
			name: "detached head",
			input: `worktree /repo/.bare
HEAD abc1234
bare

worktree /repo/detached
HEAD def5678
detached
`,
			want: []Worktree{
				{Path: "/repo/.bare", IsBare: true},
				{Path: "/repo/detached"},
			},
		},
		{
			name:  "empty output",
			input: "",
			want:  nil,
		},
		{
			name: "single normal repo",
			input: `worktree /repo
HEAD abc1234
branch refs/heads/main
`,
			want: []Worktree{
				{Path: "/repo", Branch: "main"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseWorktreeList(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d worktrees, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i].Path != tt.want[i].Path {
					t.Errorf("worktree[%d].Path = %q, want %q", i, got[i].Path, tt.want[i].Path)
				}
				if got[i].Branch != tt.want[i].Branch {
					t.Errorf("worktree[%d].Branch = %q, want %q", i, got[i].Branch, tt.want[i].Branch)
				}
				if got[i].IsBare != tt.want[i].IsBare {
					t.Errorf("worktree[%d].IsBare = %v, want %v", i, got[i].IsBare, tt.want[i].IsBare)
				}
			}
		})
	}
}
