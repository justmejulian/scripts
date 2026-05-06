package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// commentPrefix

func TestCommentPrefix(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"main.go", "//"},
		{"app.ts", "//"},
		{"component.tsx", "//"},
		{"Main.java", "//"},
		{"lib.rs", "//"},
		{"script.py", "#"},
		{"config.yaml", "#"},
		{"config.yml", "#"},
		{"Makefile.conf", "#"},
		{"query.sql", "--"},
		{"unknown.xyz", "//"},
		{"noextension", "//"},
	}
	for _, tt := range tests {
		got := commentPrefix(tt.path)
		if got != tt.want {
			t.Errorf("commentPrefix(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

// formatThreadLines

func TestFormatThreadLines_SingleComment(t *testing.T) {
	th := Thread{
		ID:       42,
		Comments: []Comment{{Author: "Alice", Content: "fix this"}},
	}
	got := formatThreadLines(th, "//")
	want := []string{"// REVU[42] @Alice: fix this"}
	assertLines(t, got, want)
}

func TestFormatThreadLines_NoPrefix(t *testing.T) {
	th := Thread{
		ID:       42,
		Comments: []Comment{{Author: "Alice", Content: "fix this"}},
	}
	got := formatThreadLines(th, "")
	want := []string{"REVU[42] @Alice: fix this"}
	assertLines(t, got, want)
}

func TestFormatThreadLines_MultiLineContent(t *testing.T) {
	th := Thread{
		ID:       42,
		Comments: []Comment{{Author: "Alice", Content: "line1\nline2\nline3"}},
	}
	got := formatThreadLines(th, "//")
	want := []string{
		"// REVU[42] @Alice: line1",
		"// REVU[42]   line2",
		"// REVU[42]   line3",
	}
	assertLines(t, got, want)
}

func TestFormatThreadLines_SkipsEmptyContentLines(t *testing.T) {
	th := Thread{
		ID:       42,
		Comments: []Comment{{Author: "Alice", Content: "line1\n\nline3"}},
	}
	got := formatThreadLines(th, "//")
	want := []string{
		"// REVU[42] @Alice: line1",
		"// REVU[42]   line3",
	}
	assertLines(t, got, want)
}

func TestFormatThreadLines_MultipleComments(t *testing.T) {
	th := Thread{
		ID: 42,
		Comments: []Comment{
			{Author: "Alice", Content: "why?"},
			{Author: "Bob", Content: "because\nreason"},
		},
	}
	got := formatThreadLines(th, "//")
	want := []string{
		"// REVU[42] @Alice: why?",
		"// REVU[42] @Bob: because",
		"// REVU[42]   reason",
	}
	assertLines(t, got, want)
}

func TestFormatThreadLines_WindowsLineEndings(t *testing.T) {
	th := Thread{
		ID:       42,
		Comments: []Comment{{Author: "Alice", Content: "line1\r\nline2\r\n"}},
	}
	got := formatThreadLines(th, "//")
	want := []string{
		"// REVU[42] @Alice: line1",
		"// REVU[42]   line2",
	}
	assertLines(t, got, want)
}

// cleanRevuLines

func TestCleanRevuLines_RemovesAllPrefixes(t *testing.T) {
	input := []string{
		"package main",
		"// REVU[1] @Alice: comment",
		"func foo() {}",
		"# REVU[2] @Bob: note",
		"-- REVU[3] @Eve: sql note",
		"return nil",
	}
	got := cleanRevuLines(input)
	want := []string{"package main", "func foo() {}", "return nil"}
	assertLines(t, got, want)
}

func TestCleanRevuLines_KeepsNonRevuLines(t *testing.T) {
	input := []string{"line1", "line2", "line3"}
	got := cleanRevuLines(input)
	assertLines(t, got, input)
}

func TestCleanRevuLines_EmptyInput(t *testing.T) {
	got := cleanRevuLines([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestCleanRevuLines_DoesNotMatchSimilarPatterns(t *testing.T) {
	input := []string{
		"// REVU without brackets",
		"// NOT_REVU[1] @Alice: not a match",
		"someREVU[1]code",
	}
	got := cleanRevuLines(input)
	assertLines(t, got, input)
}

// processFile

func TestProcessFile_InsertsAboveTargetLine(t *testing.T) {
	path := writeTempFile(t, "line1\nline2\nline3\n")

	err := processFile(path, map[int][]string{
		2: {"// REVU[1] @Alice: comment"},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := readFile(t, path)
	want := "line1\n// REVU[1] @Alice: comment\nline2\nline3\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProcessFile_MultipleInsertions(t *testing.T) {
	path := writeTempFile(t, "a\nb\nc\n")

	err := processFile(path, map[int][]string{
		1: {"// REVU[1] @Alice: on a"},
		3: {"// REVU[2] @Bob: on c"},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := readFile(t, path)
	want := "// REVU[1] @Alice: on a\na\nb\n// REVU[2] @Bob: on c\nc\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProcessFile_Idempotent(t *testing.T) {
	path := writeTempFile(t, "line1\nline2\n")

	insertions := map[int][]string{
		2: {"// REVU[1] @Alice: comment"},
	}
	processFile(path, insertions)
	processFile(path, insertions)

	got := readFile(t, path)
	if strings.Count(got, "REVU[1]") != 1 {
		t.Errorf("expected exactly 1 REVU comment after 2 syncs, got:\n%s", got)
	}
}

func TestProcessFile_CleanOnly(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[1] @Alice: old\nline2\n")

	err := processFile(path, nil)
	if err != nil {
		t.Fatal(err)
	}

	got := readFile(t, path)
	if strings.Contains(got, "REVU") {
		t.Errorf("expected REVU lines removed, got:\n%s", got)
	}
	if !strings.Contains(got, "line1") || !strings.Contains(got, "line2") {
		t.Errorf("expected original lines preserved, got:\n%s", got)
	}
}

// syncFiles

func TestSyncFiles_InjectsComments(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, "app.go")
	os.WriteFile(path, []byte("package main\n\nfunc foo() {}\n"), 0644)

	threads := []Thread{
		{
			ID:       1,
			Status:   "active",
			FilePath: "/app.go",
			Line:     3,
			Comments: []Comment{{Author: "Alice", Content: "fix this"}},
		},
	}

	if err := syncFiles(threads, repoRoot); err != nil {
		t.Fatal(err)
	}

	got := readFile(t, path)
	if !strings.Contains(got, "// REVU[1] @Alice: fix this") {
		t.Errorf("expected REVU comment injected, got:\n%s", got)
	}
}

func TestSyncFiles_ActiveOnly_SkipsResolved(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, "app.go")
	os.WriteFile(path, []byte("package main\n\nfunc foo() {}\nfunc bar() {}\n"), 0644)

	threads := []Thread{
		{
			ID:       1,
			Status:   "active",
			FilePath: "/app.go",
			Line:     3,
			Comments: []Comment{{Author: "Alice", Content: "active comment"}},
		},
		{
			ID:       2,
			Status:   "resolved",
			FilePath: "/app.go",
			Line:     4,
			Comments: []Comment{{Author: "Bob", Content: "resolved comment"}},
		},
	}

	if err := syncFiles(filterActive(threads), repoRoot); err != nil {
		t.Fatal(err)
	}

	got := readFile(t, path)
	if !strings.Contains(got, "REVU[1]") {
		t.Errorf("expected active thread injected, got:\n%s", got)
	}
	if strings.Contains(got, "REVU[2]") {
		t.Errorf("expected resolved thread skipped, got:\n%s", got)
	}
}

func TestCleanFiles_RemovesComments(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, "app.go")
	os.WriteFile(path, []byte("// REVU[1] @Alice: old\npackage main\n"), 0644)

	threads := []Thread{
		{
			ID:       1,
			Status:   "active",
			FilePath: "/app.go",
			Line:     2,
			Comments: []Comment{{Author: "Alice", Content: "new comment"}},
		},
	}

	if err := cleanFiles(threads, repoRoot); err != nil {
		t.Fatal(err)
	}

	got := readFile(t, path)
	if strings.Contains(got, "REVU") {
		t.Errorf("expected all REVU lines removed, got:\n%s", got)
	}
}

func TestSyncFiles_SkipsThreadsWithoutLocation(t *testing.T) {
	repoRoot := t.TempDir()

	threads := []Thread{
		{ID: 1, Status: "active", FilePath: "", Line: 0, Comments: []Comment{{Author: "Alice", Content: "general comment"}}},
		{ID: 2, Status: "active", FilePath: "/app.go", Line: 0, Comments: []Comment{{Author: "Bob", Content: "file but no line"}}},
	}

	// no files created — if syncFiles tries to open them it will error
	if err := syncFiles(threads, repoRoot); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// printThreads

func TestPrintThreads_WithFilePath(t *testing.T) {
	threads := []Thread{
		{
			ID:       42,
			Status:   "active",
			FilePath: "/app/foo.go",
			Line:     10,
			Comments: []Comment{{Author: "Alice", Content: "why?"}},
		},
	}

	var buf bytes.Buffer
	printThreads(&buf, threads)
	got := buf.String()

	if !strings.Contains(got, "[active] /app/foo.go:10") {
		t.Errorf("missing header, got:\n%s", got)
	}
	if !strings.Contains(got, "REVU[42] @Alice: why?") {
		t.Errorf("missing comment, got:\n%s", got)
	}
}

func TestPrintThreads_WithoutFilePath(t *testing.T) {
	threads := []Thread{
		{
			ID:       7,
			Status:   "active",
			Comments: []Comment{{Author: "Bob", Content: "general note"}},
		},
	}

	var buf bytes.Buffer
	printThreads(&buf, threads)
	got := buf.String()

	if !strings.Contains(got, "[active] thread #7") {
		t.Errorf("missing general thread header, got:\n%s", got)
	}
}

func TestPrintThreads_SeparatedByBlankLine(t *testing.T) {
	threads := []Thread{
		{ID: 1, Status: "active", FilePath: "/a.go", Line: 1, Comments: []Comment{{Author: "A", Content: "first"}}},
		{ID: 2, Status: "active", FilePath: "/b.go", Line: 2, Comments: []Comment{{Author: "B", Content: "second"}}},
	}

	var buf bytes.Buffer
	printThreads(&buf, threads)
	got := buf.String()

	if !strings.Contains(got, "\n\n") {
		t.Errorf("expected blank line between threads, got:\n%s", got)
	}
}

// cleanFile

func TestCleanFile_RemovesNumberedAndNew(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[1] @Alice: old\n// REVU[NEW] my comment\nline2\n")
	if err := cleanFile(path); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	if strings.Contains(got, "REVU") {
		t.Errorf("expected all REVU lines removed, got:\n%s", got)
	}
	if !strings.Contains(got, "line1") || !strings.Contains(got, "line2") {
		t.Errorf("expected original lines preserved, got:\n%s", got)
	}
}

// deleteRevuLine

func TestDeleteRevuLine_SingleLine(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[NEW] fix this\nline2\n")
	if err := deleteRevuLine(path, 2); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	want := "line1\nline2\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDeleteRevuLine_MultiLine(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[NEW] fix this:\n// REVU[NEW]  1. first\n// REVU[NEW]  2. second\nline2\n")
	if err := deleteRevuLine(path, 2); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	want := "line1\nline2\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// replaceNewWithID

func TestReplaceNewWithID_SingleLine(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[NEW] fix this\nline2\n")
	if err := replaceNewWithID(path, 2, 42); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	if strings.Contains(got, "REVU[NEW]") {
		t.Errorf("expected REVU[NEW] replaced, got:\n%s", got)
	}
	if !strings.Contains(got, "REVU[42]") {
		t.Errorf("expected REVU[42], got:\n%s", got)
	}
}

func TestReplaceNewWithID_MultiLine(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[NEW] fix this:\n// REVU[NEW]  1. first\n// REVU[NEW]  2. second\nline2\n")
	if err := replaceNewWithID(path, 2, 99); err != nil {
		t.Fatal(err)
	}
	got := readFile(t, path)
	if strings.Contains(got, "REVU[NEW]") {
		t.Errorf("expected all REVU[NEW] replaced, got:\n%s", got)
	}
	if strings.Count(got, "REVU[99]") != 3 {
		t.Errorf("expected 3 REVU[99] lines, got:\n%s", got)
	}
}

// readRevuNewText

func TestReadRevuNewText_SingleLine(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[NEW] fix this\nline2\n")
	text, ok := readRevuNewText(path, 2)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if text != "fix this" {
		t.Errorf("got %q, want %q", text, "fix this")
	}
}

func TestReadRevuNewText_MultiLine(t *testing.T) {
	path := writeTempFile(t, "line1\n// REVU[NEW] fix this:\n// REVU[NEW]  1. first\n// REVU[NEW]  2. second\nline2\n")
	text, ok := readRevuNewText(path, 2)
	if !ok {
		t.Fatal("expected ok=true")
	}
	want := "fix this:\n1. first\n2. second"
	if text != want {
		t.Errorf("got %q, want %q", text, want)
	}
}

// helpers

func assertLines(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %d lines, want %d\ngot:  %v\nwant: %v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d:\n got:  %q\n want: %q", i, got[i], want[i])
		}
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.go")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
