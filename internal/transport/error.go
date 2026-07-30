package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int

	Code string `json:"code"`

	Message string `json:"message"`
}

func (e *APIError) Error() string {

	return fmt.Sprintf("api error (%d): %s: %s", e.StatusCode, e.Code, e.Message)
}

func decodeError(resp *http.Response) error {

	apiErr := &APIError{
		StatusCode: resp.StatusCode,
	}

	err := json.NewDecoder(resp.Body).Decode(apiErr)

	if err != nil {

		return fmt.Errorf("http error %d", resp.StatusCode)
	}

	return apiErr
}
