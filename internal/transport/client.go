package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *HTTPClient) Do(ctx context.Context, method, path string, body, result any) error {

	var requestBody io.Reader

	if body != nil {

		data, err := json.Marshal(body)

		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}

		requestBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, requestBody)

	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("execute http request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {

		return decodeError(resp)
	}

	if result == nil {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(result)

	if err != nil {

		if err == io.EOF {
			return nil
		}

		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (c *HTTPClient) setHeaders(req *http.Request) {

	req.Header.Set("X-Auth-Token", c.Token)
	req.Header.Set("User-Agent", c.UserAgent)

	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

}
