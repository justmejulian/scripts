package main

import (
	"fmt"
	"io"
	"strings"
)

func formatThreadLines(t Thread, prefix string) []string {
	var lines []string
	for _, c := range t.Comments {
		parts := strings.Split(c.Content, "\n")
		first := true
		for _, part := range parts {
			part = strings.TrimRight(part, "\r")
			if strings.TrimSpace(part) == "" {
				continue
			}
			var line string
			if first {
				line = fmt.Sprintf("REVU[%d] @%s: %s", t.ID, c.Author, part)
				first = false
			} else {
				line = fmt.Sprintf("REVU[%d]   %s", t.ID, part)
			}
			if prefix != "" {
				line = prefix + " " + line
			}
			lines = append(lines, line)
		}
	}
	return lines
}

func printThreads(w io.Writer, threads []Thread) {
	for i, t := range threads {
		if i > 0 {
			fmt.Fprintln(w)
		}
		if t.FilePath != "" {
			fmt.Fprintf(w, "[%s] %s:%d\n", t.Status, t.FilePath, t.Line)
		} else {
			fmt.Fprintf(w, "[%s] thread #%d\n", t.Status, t.ID)
		}
		for _, line := range formatThreadLines(t, "") {
			fmt.Fprintln(w, line)
		}
	}
}
