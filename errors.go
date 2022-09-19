package dbaas

import (
	"fmt"
)

// Error titles.
const (
	ErrorNotFoundTitle   = "Not Found"
	ErrorBadRequestTitle = "Bad Request"
)

// APIError is a type of an error raised by API calls made by this library.
type APIError struct {
	APIError struct {
		Message string `json:"message"`
		Title   string `json:"title"`
		Code    int    `json:"code"`
	} `json:"error"`
}

// Error returns string representation of the error.
func (e APIError) Error() string {
	return fmt.Sprintf("%v: %v. Code: %v", e.APIError.Title, e.APIError.Message, e.APIError.Code)
}

// StatusCode returns the HTTP status from the error response.
func (e APIError) StatusCode() int {
	return e.APIError.Code
}
