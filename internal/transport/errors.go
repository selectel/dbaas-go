package transport

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrorTokenRequired    = errors.New("token is required")
	ErrorEndpointRequired = errors.New("endpoint is required")
)

// DBaaSAPIError stores useful fields from DBAAS api.
// We use a new structure to avoid breaking backward compatibility with v1 DBaaSAPIError.
type DBaaSAPIError struct {
	Method   string `json:"-"`
	Path     string `json:"-"`
	APIError struct {
		Message string `json:"message"`
	} `json:"error"`
	HTTPStatus int `json:"-"`
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

// HTTPStatusError describes interface for retry policy error.
type HTTPStatusError interface {
	error
	StatusCode() int
}
