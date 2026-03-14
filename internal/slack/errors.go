package slack

import "fmt"

type APIError struct {
	Code string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("slack: %s", e.Code)
}
