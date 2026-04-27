package main

import "context"

type Comment struct {
	ID      int
	Author  string
	Content string
}

type Thread struct {
	ID       int
	Status   string
	FilePath string
	Line     int
	Comments []Comment
}

type Provider interface {
	GetThreads(ctx context.Context, prID int) ([]Thread, error)
}

type PendingComment struct {
	AbsPath         string // absolute path to file
	RepoPath        string // path relative to repo root (e.g. /src/foo.go)
	Line            int    // 1-based line number of the REVU[NEW] line
	CodeLine        int    // line the comment is about (for Azure threadContext)
	Text            string // comment text after "REVU[NEW] "
	ReplyToThreadID int    // if > 0, add reply to this existing thread instead of creating new
}
