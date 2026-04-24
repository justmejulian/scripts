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
	Comments []Comment
}

type Provider interface {
	GetThreads(ctx context.Context, prID int) ([]Thread, error)
}
