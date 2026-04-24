package main

import (
	"context"

	"scripts/internal/azure"
)

type azureProvider struct {
	client  *azure.Client
	project string
	repo    string
}

func (p *azureProvider) GetThreads(ctx context.Context, prID int) ([]Thread, error) {
	raw, err := p.client.GetPRThreads(ctx, p.project, p.repo, prID)
	if err != nil {
		return nil, err
	}

	threads := make([]Thread, 0, len(raw))
	for _, t := range raw {
		comments := make([]Comment, 0, len(t.Comments))
		for _, c := range t.Comments {
			if c.CommentType != "text" {
				continue
			}
			comments = append(comments, Comment{
				ID:      c.ID,
				Author:  c.Author.DisplayName,
				Content: c.Content,
			})
		}
		if len(comments) == 0 {
			continue
		}
		thread := Thread{
			ID:       t.ID,
			Status:   t.Status,
			Comments: comments,
		}
		if t.ThreadContext != nil {
			thread.FilePath = t.ThreadContext.FilePath
			if t.ThreadContext.RightFileStart != nil {
				thread.Line = t.ThreadContext.RightFileStart.Line
			}
		}
		threads = append(threads, thread)
	}

	return threads, nil
}
