package azure

import "fmt"

type APIError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("azure: %s: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("azure: %s", e.Status)
}
