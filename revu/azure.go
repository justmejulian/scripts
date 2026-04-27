package main

import (
	"context"
	"fmt"

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

func (p *azureProvider) PostThread(ctx context.Context, prID int, filePath string, codeLine int, text string) (int, error) {
	req := azure.CreateThreadRequest{
		Status: "active",
		Comments: []azure.CreateThreadComment{
			{
				ParentCommentID: 0,
				Content:         text,
				CommentType:     "text",
			},
		},
		ThreadContext: &azure.ThreadContext{
			FilePath:       filePath,
			RightFileStart: &azure.FilePosition{Line: codeLine, Offset: 1},
			RightFileEnd:   &azure.FilePosition{Line: codeLine, Offset: 1},
		},
	}

	thread, err := p.client.CreatePRThread(ctx, p.project, p.repo, prID, req)
	if err != nil {
		return 0, fmt.Errorf("post thread: %w", err)
	}

	return thread.ID, nil
}

func (p *azureProvider) ReplyToThread(ctx context.Context, prID, threadID int, text string) error {
	_, err := p.client.AddThreadComment(ctx, p.project, p.repo, prID, threadID, text)
	if err != nil {
		return fmt.Errorf("reply to thread: %w", err)
	}
	return nil
}
