# prompt

Helpers for interactive terminal prompts.

## Usage

```go
import "scripts/internal/prompt"

if prompt.Confirm("push branch to remote? [y/N]: ") {
    // user answered yes
}
```

## Functions

### `Confirm(question string) bool`

Prints `question` to stderr and reads a line from stdin. Returns `true` if the user types `y` or `yes` (case-insensitive), `false` otherwise.
