package main

import (
	"errors"
	"testing"
)

func fakeGit(output string, err error) func(...string) ([]byte, error) {
	return func(args ...string) ([]byte, error) {
		return []byte(output), err
	}
}

func TestCurrentBranch(t *testing.T) {
	tests := []struct {
		name   string
		output string
		err    error
		want   string
	}{
		{"normal branch", "main\n", nil, "main"},
		{"feature branch", "feature/login\n", nil, "feature/login"},
		{"trims whitespace", "  main  \n", nil, "main"},
		{"git error returns unknown", "", errors.New("not a git repo"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{run: fakeGit(tt.output, tt.err)}
			if got := g.CurrentBranch(); got != tt.want {
				t.Errorf("CurrentBranch() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRecentLog(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{"single commit", "abc1234 feat: add thing\n", nil, "abc1234 feat: add thing", false},
		{"multiple commits", "abc1234 feat: add thing\ndef5678 fix: bug\n", nil, "abc1234 feat: add thing\ndef5678 fix: bug", false},
		{"trims whitespace", "  abc1234 feat: add thing  \n", nil, "abc1234 feat: add thing", false},
		{"git error propagates", "", errors.New("not a git repo"), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{run: fakeGit(tt.output, tt.err)}
			got, err := g.RecentLog(5)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecentLog() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("RecentLog() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStagedDiff(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		err     error
		want    string
		wantErr bool
	}{
		{"has diff", "diff --git a/main.go b/main.go\n+added line\n", nil, "diff --git a/main.go b/main.go\n+added line\n", false},
		{"empty diff", "", nil, "", false},
		{"git error propagates", "", errors.New("not a git repo"), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Git{run: fakeGit(tt.output, tt.err)}
			got, err := g.StagedDiff()
			if (err != nil) != tt.wantErr {
				t.Errorf("StagedDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("StagedDiff() = %q, want %q", got, tt.want)
			}
		})
	}
}
