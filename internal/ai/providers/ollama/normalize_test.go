package ollama

import "testing"

func TestNormalize(t *testing.T) {
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
			got := normalize(tt.input)
			if got != tt.want {
				t.Errorf("normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
