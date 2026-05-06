package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var revuLineRe = regexp.MustCompile(`^\s*\S+\s+REVU\[\d+\]`)
var revuNewLineRe = regexp.MustCompile(`^\s*\S+\s+REVU\[NEW\]`)

func commentPrefix(filePath string) string {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".go", ".js", ".ts", ".jsx", ".tsx", ".java", ".c", ".cpp", ".cs", ".swift", ".kt", ".rs":
		return "//"
	case ".py", ".rb", ".sh", ".bash", ".zsh", ".yaml", ".yml", ".toml", ".conf":
		return "#"
	case ".sql":
		return "--"
	default:
		return "//"
	}
}


func cleanRevuLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if !revuLineRe.MatchString(line) {
			result = append(result, line)
		}
	}
	return result
}

func processFile(absPath string, insertions map[int][]string) error {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	lines = cleanRevuLines(lines)

	if insertions != nil {
		lineNums := make([]int, 0, len(insertions))
		for ln := range insertions {
			lineNums = append(lineNums, ln)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(lineNums)))

		for _, ln := range lineNums {
			idx := ln - 1
			if idx < 0 {
				idx = 0
			}
			if idx > len(lines) {
				idx = len(lines)
			}
			commentLines := insertions[ln]
			lines = append(lines[:idx], append(commentLines, lines[idx:]...)...)
		}
	}

	return os.WriteFile(absPath, []byte(strings.Join(lines, "\n")), 0644)
}

func cleanFile(absPath string) error {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	lines = cleanRevuLines(lines)

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if !revuNewLineRe.MatchString(line) {
			result = append(result, line)
		}
	}

	return os.WriteFile(absPath, []byte(strings.Join(result, "\n")), 0644)
}

func filterActive(threads []Thread) []Thread {
	result := make([]Thread, 0, len(threads))
	for _, t := range threads {
		if t.Status == "active" {
			result = append(result, t)
		}
	}
	return result
}

func groupThreadsByFile(threads []Thread) map[string][]Thread {
	fileThreads := make(map[string][]Thread)
	for _, t := range threads {
		if t.FilePath == "" || t.Line == 0 {
			continue
		}
		fileThreads[t.FilePath] = append(fileThreads[t.FilePath], t)
	}
	return fileThreads
}

func syncFiles(threads []Thread, repoRoot string) error {
	for filePath, ts := range groupThreadsByFile(threads) {
		absPath := filepath.Join(repoRoot, strings.TrimPrefix(filePath, "/"))
		insertions := make(map[int][]string)
		prefix := commentPrefix(filePath)
		for _, t := range ts {
			commentLines := formatThreadLines(t, prefix)
			insertions[t.Line] = append(insertions[t.Line], commentLines...)
		}
		if err := processFile(absPath, insertions); err != nil {
			return fmt.Errorf("processing %s: %w", filePath, err)
		}
	}
	return nil
}

func cleanFiles(threads []Thread, repoRoot string) error {
	for filePath := range groupThreadsByFile(threads) {
		absPath := filepath.Join(repoRoot, strings.TrimPrefix(filePath, "/"))
		if err := cleanFile(absPath); err != nil {
			return fmt.Errorf("cleaning %s: %w", filePath, err)
		}
	}
	return nil
}
