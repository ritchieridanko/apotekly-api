package caches

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const CacheErrorTracer = ce.CacheTracer

type Cache interface {
	Del(ctx context.Context, keys ...string) (err error)
	ShouldAccountBeLocked(ctx context.Context, key string) (shouldBeLocked bool, err error)
}

type cache struct {
	client     *redis.Client
	maxRetries int
	baseDelay  int
}

func NewCache(client *redis.Client) Cache {
	return &cache{client, config.GetCacheMaxRetries(), config.GetCacheBaseDelay()}
}
