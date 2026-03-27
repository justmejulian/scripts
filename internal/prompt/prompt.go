package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks the user a yes/no question and returns true if they answer "y" or "yes".
func Confirm(question string) bool {
	fmt.Fprint(os.Stderr, question)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		return answer == "y" || answer == "yes"
	}
	return false
}
