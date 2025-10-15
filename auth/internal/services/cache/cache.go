package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	cache *redis.Client

	maxRetries int
	baseDelay  int
}

func NewCache(cache *redis.Client, maxRetries, baseDelay int) *Cache {
	return &Cache{cache, maxRetries, baseDelay}
}

func (c *Cache) Evaluate(ctx context.Context, hashKey, script string, keys []string, args ...interface{}) (interface{}, error) {
	hash, err := c.Get(ctx, hashKey)
	if err != nil {
		hash, err = c.loadScript(ctx, script)
		if err != nil {
			return nil, err
		}

		// set indefinitely
		if err := c.Set(ctx, hashKey, hash, -1); err != nil {
			return nil, err
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		result, err := c.cache.EvalSha(ctx, hash, keys, args...).Result()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return nil, err
		}
	}

	return nil, lastErr
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		var err error
		if duration <= 0 {
			err = c.cache.Set(ctx, key, value, 0).Err()
		} else {
			err = c.cache.Set(ctx, key, value, duration).Err()
		}

		if err == nil {
			return nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return err
		}
	}

	return lastErr
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		result, err := c.cache.Get(ctx, key).Result()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return "", err
		}
	}

	return "", lastErr
}

func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		err := c.cache.Del(ctx, keys...).Err()
		if err == nil {
			return nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return err
		}
	}

	return lastErr
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		result, err := c.cache.Exists(ctx, key).Result()
		if err == nil {
			return result > 0, nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return false, err
		}
	}

	return false, lastErr
}

func (c *Cache) loadScript(ctx context.Context, script string) (string, error) {
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		result, err := c.cache.ScriptLoad(ctx, script).Result()
		if err == nil {
			return result, nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.baseDelay, attempt); err != nil {
			return "", err
		}
	}

	return "", lastErr
}
