package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git struct {
	run func(args ...string) ([]byte, error)
}

func NewGit() *Git {
	return &Git{
		run: func(args ...string) ([]byte, error) {
			return exec.Command("git", args...).Output()
		},
	}
}

func (g *Git) CurrentBranch() string {
	out, err := g.run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func (g *Git) RecentLog(n int) (string, error) {
	out, err := g.run("log", fmt.Sprintf("-n%d", n), "--oneline")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (g *Git) StagedDiff() (string, error) {
	out, err := g.run("diff", "--cached")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
