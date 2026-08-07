package transport

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxRetries int

	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

type RetryPolicy interface {
	ShouldRetry(err error, attempt int) bool
	Delay(attempt int) time.Duration
}

type defaultRetryPolicy struct {
	config RetryConfig
}

func (r *defaultRetryPolicy) ShouldRetry(err error, attempt int) bool {
	if attempt >= r.config.MaxRetries {
		return false
	}

	var httpErr HTTPStatusError

	if errors.As(err, &httpErr) {
		switch httpErr.StatusCode() {
		case
			http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:

			return true
		}

		return false
	}

	var netErr net.Error

	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return errors.Is(err, context.DeadlineExceeded)
}

// Delay calculates exponential backoff delay.
func (r *defaultRetryPolicy) Delay(attempt int) time.Duration {
	delay := r.config.InitialBackoff * time.Duration(1<<attempt)

	if delay > r.config.MaxBackoff {
		return r.config.MaxBackoff
	}
	return delay
}

func NewRetry(config RetryConfig) RetryPolicy {
	if config.InitialBackoff <= 0 {
		config.InitialBackoff = time.Second
	}

	if config.MaxBackoff <= 0 {
		config.MaxBackoff = 30 * time.Second
	}

	return &defaultRetryPolicy{
		config: config,
	}
}
