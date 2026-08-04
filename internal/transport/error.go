package transport

import (
	"fmt"
)

type APIError struct {
	StatusCode int

	Code string `json:"code"`

	Message string `json:"message"`
}

func (e *APIError) Error() string {

	return fmt.Sprintf("api error (%d): %s: %s", e.StatusCode, e.Code, e.Message)
}

// DBaaSAPIError is a type of an error raised by API calls made by this library.
type DBaaSAPIError struct {
	APIError struct {
		Message string `json:"message"`
		Title   string `json:"title"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// Error returns string representation of the error.
func (e DBaaSAPIError) Error() string {
	return fmt.Sprintf("%v: %v. Code: %v", e.APIError.Title, e.APIError.Message, e.APIError.Code)
}

// StatusCode returns the HTTP status from the error response.
func (e DBaaSAPIError) StatusCode() int {
	return e.APIError.Code
}
