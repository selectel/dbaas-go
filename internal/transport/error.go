package transport

import (
	"encoding/json"
	"fmt"
)

type DBaaSAPIError struct {
	StatusCode int    `json:"-"`
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
		e.StatusCode,
		e.APIError.Message,
	)
}

func decodeError(statusCode int, method string, path string, body []byte) error {

	apiErr := &DBaaSAPIError{
		StatusCode: statusCode,
		Method:     method,
		Path:       path,
	}

	if err := json.Unmarshal(body, apiErr); err != nil {
		return fmt.Errorf("could not unmarshal error response: %w", err)
	}

	return apiErr
}
