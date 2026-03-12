package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func currentBranch() string {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func recentLog(n int) (string, error) {
	out, err := exec.Command("git", "log", fmt.Sprintf("-n%d", n), "--oneline").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func stagedDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
