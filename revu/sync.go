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

func processFile(absPath string, insertions map[int][]string, cleanNew bool) error {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	lines = cleanRevuLines(lines)
	if cleanNew {
		result := make([]string, 0, len(lines))
		for _, line := range lines {
			if !revuNewLineRe.MatchString(line) {
				result = append(result, line)
			}
		}
		lines = result
	}

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

func syncFiles(threads []Thread, repoRoot string, cleanOnly bool, activeOnly bool) error {
	fileThreads := make(map[string][]Thread)
	for _, t := range threads {
		if t.FilePath == "" || t.Line == 0 {
			continue
		}
		if activeOnly && t.Status != "active" {
			continue
		}
		fileThreads[t.FilePath] = append(fileThreads[t.FilePath], t)
	}

	for filePath, ts := range fileThreads {
		absPath := filepath.Join(repoRoot, strings.TrimPrefix(filePath, "/"))

		var insertions map[int][]string
		if !cleanOnly {
			insertions = make(map[int][]string)
			prefix := commentPrefix(filePath)
			for _, t := range ts {
				commentLines := formatThreadLines(t, prefix)
				insertions[t.Line] = append(insertions[t.Line], commentLines...)
			}
		}

		if err := processFile(absPath, insertions, cleanOnly); err != nil {
			return fmt.Errorf("processing %s: %w", filePath, err)
		}
	}

	return nil
}
