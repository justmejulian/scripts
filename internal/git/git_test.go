package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCurrentBranch_Success(t *testing.T) {
	dir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test")
	// Create an initial commit so we can switch branches
	f := filepath.Join(dir, "README.md")
	if err := os.WriteFile(f, []byte("hi"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", ".")
	run("git", "commit", "-m", "init")
	run("git", "checkout", "-b", "feat/PROJ-123-test")

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	branch, err := CurrentBranch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "feat/PROJ-123-test" {
		t.Errorf("got %q, want %q", branch, "feat/PROJ-123-test")
	}
}

func TestCurrentBranch_NotARepo(t *testing.T) {
	dir := t.TempDir()

	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(orig)

	_, err = CurrentBranch()
	if err == nil {
		t.Fatal("expected error when not in a git repo")
	}
}
