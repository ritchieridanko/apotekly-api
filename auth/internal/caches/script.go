package caches

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

func (c *cache) ShouldAccountBeLocked(ctx context.Context, key string) (bool, error) {
	tracer := CacheErrorTracer + ": ShouldAccountBeLocked()"

	script := redis.NewScript(`
		local count = redis.call("INCR", KEYS[1])
		if count == 1 then
			redis.call("EXPIRE", KEYS[1], ARGV[2])
		end
		if count >= tonumber(ARGV[1]) then
			return 1
		end
		return 0
	`)

	result, err := script.Run(ctx, c.client, []string{key}, constants.MaxTotalFailedAuth, constants.RedisDurationTotalFailedAuth).Int()
	if err != nil {
		return false, ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
	}

	return result == 1, nil
}

func (c *cache) NewOrReplacePasswordResetToken(ctx context.Context, authID int64, token string, duration time.Duration) error {
	tracer := CacheErrorTracer + ": NewOrReplacePasswordResetToken()"

	authKey := utils.GenerateDynamicRedisKey(constants.RedisKeyPasswordResetAuth, authID)
	tokenKey := utils.GenerateDynamicRedisKey(constants.RedisKeyPasswordResetToken, token)

	script := redis.NewScript(`
		local oldToken = redis.call("GET", KEYS[1])
		if oldToken then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. oldToken)
		end
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`)

	_, err := script.Run(ctx, c.client, []string{authKey, tokenKey, constants.RedisKeyPasswordResetToken}, token, authID, int(duration.Seconds())).Result()
	if err != nil {
		return ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
	}

	return nil
}
