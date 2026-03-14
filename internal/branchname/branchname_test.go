package branchname

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
			got := SlugifyDescription(tt.input)
			if got != tt.want {
				t.Errorf("SlugifyDescription(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBranchType(t *testing.T) {
	tests := []struct {
		branch  string
		want    string
		wantErr bool
	}{
		{"feat/PROJ-123-add-login", "feat", false},
		{"fix/PROJ-42-some-bug", "fix", false},
		{"chore/update-deps", "chore", false},
		{"no-slash-here", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			got, err := BranchType(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Fatalf("BranchType(%q) error = %v, wantErr %v", tt.branch, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("BranchType(%q) = %q, want %q", tt.branch, got, tt.want)
			}
		})
	}
}

func TestBuildName(t *testing.T) {
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
			got := BuildName(tt.taskType, tt.issue, tt.description)
			if got != tt.want {
				t.Errorf("BuildName(%q, %q, %q) = %q, want %q",
					tt.taskType, tt.issue, tt.description, got, tt.want)
			}
		})
	}
}
