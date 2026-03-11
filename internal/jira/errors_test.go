package jira

import "testing"

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Status: "404 Not Found"}
	if err.Error() != "jira: 404 Not Found" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}

func TestAPIError_Error_500(t *testing.T) {
	err := &APIError{StatusCode: 500, Status: "500 Internal Server Error"}
	if err.Error() != "jira: 500 Internal Server Error" {
		t.Errorf("unexpected error string: %s", err.Error())
	}
}
