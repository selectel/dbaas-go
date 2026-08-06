package transport

// Option is extension for HTTPClient.
type Option func(*HTTPClient)

// WithRetry provides retry option for client from retry config.
func WithRetry(config RetryConfig) Option {
	return func(c *HTTPClient) {
		c.retry = NewRetry(config)
	}
}
