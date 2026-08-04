package transport

type Option func(*HTTPClient)

func WithRetry(
	config RetryConfig,
) Option {
	return func(c *HTTPClient) {
		c.retry = NewRetry(config)
	}
}
