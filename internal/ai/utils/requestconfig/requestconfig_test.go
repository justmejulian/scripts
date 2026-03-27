package requestconfig

import (
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		input    string
		want     map[string]any
		wantErr  string
	}{
		{
			name:     "empty config",
			provider: "zen",
			input:    "",
			want:     nil,
		},
		{
			name:     "whitespace config",
			provider: "zen",
			input:    "  \n\t  ",
			want:     nil,
		},
		{
			name:     "valid object",
			provider: "ollama",
			input:    `{"temperature":0.2,"nested":{"enabled":true}}`,
			want: map[string]any{
				"temperature": 0.2,
				"nested": map[string]any{
					"enabled": true,
				},
			},
		},
		{
			name:     "invalid json",
			provider: "zen",
			input:    `{"temperature":`,
			wantErr:  `zen: invalid request config:`,
		},
		{
			name:     "non object json array",
			provider: "zen",
			input:    `[]`,
			wantErr:  `zen: invalid request config: must be a JSON object`,
		},
		{
			name:     "non object json null",
			provider: "zen",
			input:    `null`,
			wantErr:  `zen: invalid request config: must be a JSON object`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJSON(tt.provider, tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("ParseJSON(%q, %q) error = nil, want %q", tt.provider, tt.input, tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("ParseJSON(%q, %q) error = %q, want substring %q", tt.provider, tt.input, err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseJSON(%q, %q) error = %v", tt.provider, tt.input, err)
			}

			if !mapsEqual(got, tt.want) {
				t.Fatalf("ParseJSON(%q, %q) = %#v, want %#v", tt.provider, tt.input, got, tt.want)
			}
		})
	}
}

func TestEnsureNoOverrides(t *testing.T) {
	err := EnsureNoOverrides("ollama", map[string]any{"stream": false}, "model", "messages", "stream")
	if err == nil {
		t.Fatal("EnsureNoOverrides() error = nil, want error")
	}

	if got, want := err.Error(), `ollama: request config may not override "stream"`; got != want {
		t.Fatalf("EnsureNoOverrides() error = %q, want %q", got, want)
	}
}

func TestMerge(t *testing.T) {
	base := map[string]any{
		"model":     "qwen3:8b",
		"stream":    false,
		"unchanged": true,
	}
	extra := map[string]any{
		"stream":      true,
		"temperature": 0.2,
	}

	got := Merge(base, extra)

	want := map[string]any{
		"model":       "qwen3:8b",
		"stream":      true,
		"unchanged":   true,
		"temperature": 0.2,
	}

	if !mapsEqual(got, want) {
		t.Fatalf("Merge() = %#v, want %#v", got, want)
	}

	if base["stream"] != false {
		t.Fatalf("Merge() mutated base map: %#v", base)
	}
}

func TestApply(t *testing.T) {
	base := map[string]any{
		"model": "glm-5",
	}

	got, err := Apply("zen", base, `{"thinking":{"type":"disabled"}}`, "model", "messages")
	if err != nil {
		t.Fatalf("Apply() error = %v", err)
	}

	want := map[string]any{
		"model": "glm-5",
		"thinking": map[string]any{
			"type": "disabled",
		},
	}

	if !mapsEqual(got, want) {
		t.Fatalf("Apply() = %#v, want %#v", got, want)
	}

	if !mapsEqual(base, map[string]any{"model": "glm-5"}) {
		t.Fatalf("Apply() mutated base map: %#v", base)
	}
}

func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}

	for key, aValue := range a {
		bValue, ok := b[key]
		if !ok {
			return false
		}

		if !valuesEqual(aValue, bValue) {
			return false
		}
	}

	return true
}

func valuesEqual(a, b any) bool {
	aMap, aIsMap := a.(map[string]any)
	bMap, bIsMap := b.(map[string]any)
	if aIsMap || bIsMap {
		return aIsMap && bIsMap && mapsEqual(aMap, bMap)
	}

	return a == b
}
