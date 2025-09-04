package caches

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const CacheErrorTracer = ce.CacheTracer

type Cache interface {
	Has(ctx context.Context, key string) (exists bool, err error)
	Del(ctx context.Context, keys ...string) (err error)
	ShouldAccountBeLocked(ctx context.Context, key string) (shouldBeLocked bool, err error)
	NewOrReplacePasswordResetToken(ctx context.Context, authID int64, token string, duration time.Duration) (err error)
	ConsumePasswordResetToken(ctx context.Context, token string) (authID int64, err error)
	NewOrReplaceVerificationToken(ctx context.Context, authID int64, token string, duration time.Duration) (err error)
}

type cache struct {
	client     *redis.Client
	maxRetries int
	baseDelay  int
}

func NewCache(client *redis.Client) Cache {
	return &cache{client, config.GetRedisMaxRetries(), config.GetRedisBaseDelay()}
}
