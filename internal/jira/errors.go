package jira

import "fmt"

type APIError struct {
	StatusCode int
	Status     string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("jira: %s", e.Status)
}
