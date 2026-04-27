package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var newRevuRe = regexp.MustCompile(`^(\s*\S+)\s+REVU\[NEW\]\s+(.+)$`)
var existingRevuRe = regexp.MustCompile(`^\s*\S+\s+REVU\[(\d+)\]`)

func scanNewComments(repoRoot string) ([]PendingComment, error) {
	out, err := exec.Command("git", "-C", repoRoot, "ls-files", "--cached").Output()
	if err != nil {
		return nil, fmt.Errorf("git ls-files: %w", err)
	}

	var results []PendingComment
	paths := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, rel := range paths {
		if rel == "" {
			continue
		}
		abs := filepath.Join(repoRoot, rel)
		data, err := os.ReadFile(abs)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			m := newRevuRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			lineNum := i + 1
			codeLineNum := lineNum + 1
			if codeLineNum > len(lines) {
				codeLineNum = lineNum
			}
			pc := PendingComment{
				AbsPath:  abs,
				RepoPath: "/" + filepath.ToSlash(rel),
				Line:     lineNum,
				CodeLine: codeLineNum,
				Text:     strings.TrimSpace(m[2]),
			}
			// If the line immediately above belongs to an existing REVU thread,
			// this comment is a reply to that thread.
			if i > 0 {
				if tm := existingRevuRe.FindStringSubmatch(lines[i-1]); tm != nil {
					if id, err := strconv.Atoi(tm[1]); err == nil {
						pc.ReplyToThreadID = id
					}
				}
			}
			results = append(results, pc)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].AbsPath != results[j].AbsPath {
			return results[i].AbsPath < results[j].AbsPath
		}
		return results[i].Line < results[j].Line
	})

	return results, nil
}

func renderContext(absPath string, revuLine, contextLines int) string {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")

	start := revuLine - 1 - contextLines
	if start < 0 {
		start = 0
	}
	end := revuLine - 1 + contextLines
	if end >= len(lines) {
		end = len(lines) - 1
	}

	var sb strings.Builder
	for i := start; i <= end; i++ {
		lineNum := i + 1
		if lineNum == revuLine {
			fmt.Fprintf(&sb, "->%d: %s\n", lineNum, lines[i])
		} else {
			fmt.Fprintf(&sb, "  %d: %s\n", lineNum, lines[i])
		}
	}
	return sb.String()
}

func interactiveReview(comments []PendingComment, in io.Reader, out io.Writer) ([]PendingComment, error) {
	scanner := bufio.NewScanner(in)
	total := len(comments)
	var approved []PendingComment
	var toDelete []PendingComment

	for idx, c := range comments {
		for {
			if c.ReplyToThreadID > 0 {
				fmt.Fprintf(out, "\n--- (%d/%d) %s:%d [reply to thread #%d] ---\n\n", idx+1, total, c.RepoPath, c.Line, c.ReplyToThreadID)
			} else {
				fmt.Fprintf(out, "\n--- (%d/%d) %s:%d ---\n\n", idx+1, total, c.RepoPath, c.Line)
			}
			fmt.Fprint(out, renderContext(c.AbsPath, c.Line, 4))
			fmt.Fprintf(out, "\nUpload? [y]es / [n]o / [d]elete / [e]dit: ")

			if !scanner.Scan() {
				return approved, scanner.Err()
			}
			choice := strings.TrimSpace(strings.ToLower(scanner.Text()))

			switch choice {
			case "y":
				approved = append(approved, c)
				goto next
			case "n":
				goto next
			case "d":
				toDelete = append(toDelete, c)
				goto next
			case "e":
				edited, err := editInEditor(c.Text)
				if err != nil {
					fmt.Fprintf(out, "editor error: %v\n", err)
					continue
				}
				c.Text = edited
			default:
				fmt.Fprintln(out, "invalid choice")
			}
		}
	next:
	}

	// Delete in reverse order so line numbers stay valid within each file.
	sort.Slice(toDelete, func(i, j int) bool {
		if toDelete[i].AbsPath != toDelete[j].AbsPath {
			return toDelete[i].AbsPath < toDelete[j].AbsPath
		}
		return toDelete[i].Line > toDelete[j].Line
	})
	for _, c := range toDelete {
		if err := deleteRevuLine(c.AbsPath, c.Line); err != nil {
			fmt.Fprintf(out, "delete %s:%d: %v\n", c.RepoPath, c.Line, err)
		}
	}

	return approved, nil
}

func editInEditor(text string) (string, error) {
	tmp, err := os.CreateTemp("", "revu-*.txt")
	if err != nil {
		return text, err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(text); err != nil {
		tmp.Close()
		return text, err
	}
	tmp.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return text, err
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		return text, err
	}
	return strings.TrimSpace(string(data)), nil
}

func deleteRevuLine(absPath string, lineNum int) error {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("line %d out of range", lineNum)
	}
	lines = append(lines[:lineNum-1], lines[lineNum:]...)
	return os.WriteFile(absPath, []byte(strings.Join(lines, "\n")), 0)
}

func applyUploads(ctx context.Context, provider *azureProvider, prID int, approved []PendingComment, out io.Writer) {
	fmt.Fprintf(out, "\nUploading %d comment(s)...\n", len(approved))
	for idx, c := range approved {
		var threadID int
		var err error

		if c.ReplyToThreadID > 0 {
			err = provider.ReplyToThread(ctx, prID, c.ReplyToThreadID, c.Text)
			threadID = c.ReplyToThreadID
		} else {
			threadID, err = provider.PostThread(ctx, prID, c.RepoPath, c.CodeLine, c.Text)
		}

		if err != nil {
			fmt.Fprintf(out, "  [%d/%d] %s:%d -> ERROR: %v\n", idx+1, len(approved), c.RepoPath, c.Line, err)
			continue
		}
		if err := replaceNewWithID(c.AbsPath, c.Line, threadID); err != nil {
			fmt.Fprintf(out, "  [%d/%d] %s:%d -> thread #%d OK (file update failed: %v)\n", idx+1, len(approved), c.RepoPath, c.Line, threadID, err)
			continue
		}
		fmt.Fprintf(out, "  [%d/%d] %s:%d -> thread #%d OK\n", idx+1, len(approved), c.RepoPath, c.Line, threadID)
	}
	fmt.Fprintln(out, "Done. Run `revu sync` to pull uploaded comments back with full author/content formatting.")
}

func replaceNewWithID(absPath string, lineNum, threadID int) error {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if lineNum < 1 || lineNum > len(lines) {
		return fmt.Errorf("line %d out of range", lineNum)
	}
	lines[lineNum-1] = strings.Replace(lines[lineNum-1], "REVU[NEW]", fmt.Sprintf("REVU[%d]", threadID), 1)
	return os.WriteFile(absPath, []byte(strings.Join(lines, "\n")), 0)
}
