package fzf

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// ErrCancelled is returned when the user exits fzf without making a selection.
var ErrCancelled = errors.New("selection cancelled")

// Select runs fzf with the given items and returns the chosen line.
// prompt is shown as the fzf prompt (e.g. "Task: ").
func Select(prompt string, items []string) (string, error) {
	cmd := exec.Command("fzf", "--prompt", prompt+" ")

	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 130 {
			return "", ErrCancelled
		}
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
