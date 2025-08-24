package caches

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

func (c *cache) Del(ctx context.Context, keys ...string) error {
	tracer := CacheErrorTracer + ": Del()"

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		err := c.client.Del(ctx, keys...).Err()
		if err == nil {
			return nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}

		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
		}
	}

	return ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, lastErr)
}

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
