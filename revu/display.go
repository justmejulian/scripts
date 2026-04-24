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
		fmt.Fprintf(w, "[%s] thread #%d\n", t.Status, t.ID)
		for j, c := range t.Comments {
			if j == 0 {
				fmt.Fprintf(w, "  %s: %s\n", c.Author, c.Content)
			} else {
				fmt.Fprintf(w, "    %s: %s\n", c.Author, c.Content)
			}
		}
	}
}
