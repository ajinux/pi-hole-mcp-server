package client

import "fmt"

// httpError represents an HTTP error response with status code and body
type httpError struct {
	StatusCode int
	Status     string
	Body       string
	Method     string
	URL        string
}

func (e *httpError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("HTTP %d %s: %s", e.StatusCode, e.Status, e.Body)
	}
	return fmt.Sprintf("HTTP %d %s", e.StatusCode, e.Status)
}

// NewHTTPError creates a new httpError
func newHTTPError(statusCode int, status, body, method, url string) *httpError {
	return &httpError{
		StatusCode: statusCode,
		Status:     status,
		Body:       body,
		Method:     method,
		URL:        url,
	}
}
