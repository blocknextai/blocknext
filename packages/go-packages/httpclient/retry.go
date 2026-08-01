package httpclient

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

func (r *RequestBuilder) shouldRetry(err error, statusCode int) bool {
	if err != nil {
		if netErr, ok := errors.AsType[net.Error](err); ok && netErr.Timeout() {
			return true
		}

		var dnsErr *net.DNSError
		return errors.As(err, &dnsErr)
	}

	return statusCode >= http.StatusInternalServerError || statusCode == http.StatusTooManyRequests
}

func (r *RequestBuilder) executeWithRetry(ctx context.Context, fn func() (*Response, error)) (*Response, error) {
	maxRetries := r.client.retryConfig.MaxRetries
	backoffMs := r.client.retryConfig.BackoffMs

	var lastErr error
	var lastResp *Response

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return lastResp, err
		}

		resp, err := fn()

		if err == nil && !r.shouldRetry(nil, resp.Status) {
			return resp, nil
		}

		lastErr = err
		lastResp = resp

		if attempt == maxRetries {
			break
		}

		shouldRetry := r.shouldRetry(err, 0)
		if resp != nil {
			shouldRetry = shouldRetry || r.shouldRetry(nil, resp.Status)
		}

		if !shouldRetry {
			break
		}

		backoffDuration := time.Duration(backoffMs*(attempt+1)) * time.Millisecond

		select {
		case <-ctx.Done():
			return lastResp, ctx.Err()
		case <-time.After(backoffDuration):
		}
	}

	if lastErr != nil {
		return lastResp, lastErr
	}

	return lastResp, ErrMaxRetriesExceeded
}
