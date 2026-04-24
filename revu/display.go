package main

import (
	"fmt"
	"io"
)

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
		for _, c := range t.Comments {
			fmt.Fprintf(w, "%s:\n%s\n", c.Author, c.Content)
		}
	}
}
