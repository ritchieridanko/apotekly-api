package cache

import (
	"context"
	"errors"
	"time"
)

func isRetryable(err error) (isRetryable bool) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}
	return true
}

func backoffWait(ctx context.Context, baseDelay, attempt int) (err error) {
	backoff := time.Duration(baseDelay) * (1 << attempt)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(backoff):
		return nil
	}
}
