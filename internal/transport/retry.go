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
	Backoff    time.Duration
}

func isRetryable(err error) bool {
	var apiErr *DBaaSAPIError

	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusTooManyRequests,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout:
			return true
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}

	return errors.Is(err, context.DeadlineExceeded)
}

func sleepBackoff(ctx context.Context, attempt int, base time.Duration) error {

	delay := base * time.Duration(1<<attempt)

	timer := time.NewTimer(delay)

	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-timer.C:
		return nil
	}
}
