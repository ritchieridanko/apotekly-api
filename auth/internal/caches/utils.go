package caches

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

func (c *cache) Has(ctx context.Context, key string) (bool, error) {
	tracer := CacheErrorTracer + ": Has()"

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		result, err := c.client.Exists(ctx, key).Result()
		if err == nil {
			return result > 0, nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}

		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return false, ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
		}
	}

	return false, ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, lastErr)
}

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
