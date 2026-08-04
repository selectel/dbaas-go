package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	Do(ctx context.Context, method, path string, body, result any) error
}

type HTTPClient struct {
	HTTPClient *http.Client
	Token      string
	Endpoint   string
	UserAgent  string
}

func NewHTTPClient(httpClient *http.Client, token, endpoint, userAgent string) *HTTPClient {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &HTTPClient{
		HTTPClient: httpClient,
		Token:      token,
		Endpoint:   endpoint,
		UserAgent:  userAgent,
	}
}

func marshalBody(body any) (io.Reader, error) {

	if body == nil {
		return nil, nil
	}

	if b, ok := body.([]byte); ok {
		return bytes.NewReader(b), nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling params to JSON, %w", err)
	}

	return bytes.NewReader(data), nil
}

func (c *HTTPClient) Do(ctx context.Context, method, path string, body, result any) error {

	requestBody, err := marshalBody(body)
	if err != nil {
		return err
	}

	request, err := c.newRequest(ctx, method, path, requestBody)
	if err != nil {
		return err
	}

	response, err := c.HTTPClient.Do(request)

	if err != nil {
		return fmt.Errorf("HTTP request failed, %w", err)
	}

	defer response.Body.Close()

	return decodeResponse(request, response, result)
}

func (c *HTTPClient) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {

	fullURL, err := url.JoinPath(c.Endpoint, path)

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("X-Auth-Token", c.Token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func decodeResponse(req *http.Request, resp *http.Response, result any) error {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body, %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return HandleStatusCode(resp.StatusCode, body, req.URL.String())
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if result == nil {
		return nil
	}

	if len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("could not decode response body: %w", err)
	}

	return nil
}

// HandleStatusCode checks status code and returns corresponding error.
func HandleStatusCode(statusCode int, body []byte, path string) error {
	if statusCode >= http.StatusInternalServerError {
		return fmt.Errorf("http status %d: service failed.\n%v\n%v", statusCode, body, path) //nolint
	}

	errBody := &DBaaSAPIError{}
	err := json.Unmarshal(body, &errBody)
	if err != nil {
		return fmt.Errorf("can't unmarshal response:\n%s, %w", body, err)
	}
	return errBody
}
