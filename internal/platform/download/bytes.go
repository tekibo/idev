package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// Bytes fetches a small remote document with retry.
func Bytes(ctx context.Context, url string) ([]byte, error) {
	var last error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		data, retry, err := bytesOnce(ctx, url)
		if err == nil {
			return data, nil
		}
		last = err
		if !retry {
			return nil, err
		}
		if attempt < maxAttempts {
			wait(ctx, attempt)
		}
	}
	return nil, fmt.Errorf("download: %s: %w", url, last)
}

func bytesOnce(ctx context.Context, url string) ([]byte, bool, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, retryable(err), err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, retryStatus(response.StatusCode),
			fmt.Errorf("download status %s", response.Status)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		return nil, retryable(err), err
	}
	return data, false, nil
}

func wait(ctx context.Context, attempt int) {
	delay := time.Second << uint(min(attempt-1, 4))
	if delay > maxDelay {
		delay = maxDelay
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

func retryable(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	var netErr net.Error
	return errors.As(err, &netErr) || errors.Is(err, io.ErrUnexpectedEOF)
}

func retryStatus(code int) bool {
	return code == http.StatusRequestTimeout ||
		code == http.StatusTooManyRequests ||
		code >= http.StatusInternalServerError
}
