package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	Do(ctx context.Context, method, path string, body, result any) error
}

type HTTPClient struct {
	retry      RetryPolicy
	httpClient *http.Client
	endpoint   string
	userAgent  string
	token      string
}

func NewHTTPClient(httpClient *http.Client, token, endpoint, userAgent string, options ...Option) *HTTPClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &HTTPClient{
		httpClient: httpClient,
		token:      token,
		endpoint:   endpoint,
		userAgent:  userAgent,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *HTTPClient) Do(ctx context.Context, method, path string, body, result any) error {
	for attempt := 0; ; attempt++ {
		err := c.doOnce(ctx, method, path, body, result)

		if err == nil {
			return nil
		}

		if c.retry == nil {
			return err
		}

		if !c.retry.ShouldRetry(err, attempt) {
			return err
		}

		if err := wait(ctx, c.retry.Delay(attempt)); err != nil {
			return err
		}
	}
}

func (c *HTTPClient) doOnce(ctx context.Context, method, path string, body, result any) error {
	requestBody, err := marshalBody(body)
	if err != nil {
		return err
	}

	request, err := c.newRequest(ctx, method, path, requestBody)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("HTTP request failed, %w", err)
	}

	defer response.Body.Close()

	return decodeResponse(request, response, result)
}

func (c *HTTPClient) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	fullURL, err := url.JoinPath(c.endpoint, path)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP request creation failed, %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("X-Auth-Token", c.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
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

func decodeResponse(req *http.Request, resp *http.Response, result any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body, %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return decodeError(resp.StatusCode, req.Method, req.URL.Path, body)
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

func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		//nolint:wrapcheck
		return ctx.Err()

	case <-timer.C:
		return nil
	}
}
