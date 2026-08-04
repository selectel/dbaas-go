package transport

import (
	"encoding/json"
	"fmt"
)

type DBaaSAPIError struct {
	HTTPStatus int    `json:"-"`
	Method     string `json:"-"`
	Path       string `json:"-"`

	APIError struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *DBaaSAPIError) Error() string {
	return fmt.Sprintf(
		"%s %s: http %d: %s",
		e.Method,
		e.Path,
		e.HTTPStatus,
		e.APIError.Message,
	)
}

func (e *DBaaSAPIError) StatusCode() int {
	return e.HTTPStatus
}

func decodeError(statusCode int, method string, path string, body []byte) error {

	apiErr := &DBaaSAPIError{
		HTTPStatus: statusCode,
		Method:     method,
		Path:       path,
	}

	if err := json.Unmarshal(body, apiErr); err != nil {
		return fmt.Errorf("could not unmarshal error response: %w", err)
	}

	return apiErr
}

// HTTPStatusError describes interface for retry policy
type HTTPStatusError interface {
	error
	StatusCode() int
}
