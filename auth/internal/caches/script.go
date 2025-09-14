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

func (c *cache) ConsumePasswordResetToken(ctx context.Context, token string) (int64, error) {
	tracer := CacheErrorTracer + ": ConsumePasswordResetToken()"

	tokenKey := utils.GenerateDynamicRedisKey(constants.RedisKeyPasswordResetToken, token)

	script := redis.NewScript(`
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`)

	result, err := script.Run(ctx, c.client, []string{tokenKey, constants.RedisKeyPasswordResetAuth}).Result()
	if err == redis.Nil {
		return 0, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidToken, tracer, err)
	}
	if err != nil {
		return 0, ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
	}

	value, ok := result.(string)
	if !ok {
		return 0, ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, ce.ErrInvalidType)
	}

	authID, err := utils.ToInt64(value)
	if err != nil {
		return 0, ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
	}

	return authID, nil
}

func (c *cache) NewOrReplaceVerificationToken(ctx context.Context, authID int64, token string, duration time.Duration) error {
	tracer := CacheErrorTracer + ": NewOrReplaceVerificationToken()"

	authKey := utils.GenerateDynamicRedisKey(constants.RedisKeyVerificationAuth, authID)
	tokenKey := utils.GenerateDynamicRedisKey(constants.RedisKeyVerificationToken, token)

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

	_, err := script.Run(ctx, c.client, []string{authKey, tokenKey, constants.RedisKeyVerificationToken}, token, authID, int(duration.Seconds())).Result()
	if err != nil {
		return ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
	}

	return nil
}

func (c *cache) ConsumeVerificationToken(ctx context.Context, token string) (int64, error) {
	tracer := CacheErrorTracer + ": ConsumeVerificationToken()"

	tokenKey := utils.GenerateDynamicRedisKey(constants.RedisKeyVerificationToken, token)

	script := redis.NewScript(`
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`)

	result, err := script.Run(ctx, c.client, []string{tokenKey, constants.RedisKeyVerificationAuth}).Result()
	if err == redis.Nil {
		return 0, ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidToken, tracer, err)
	}
	if err != nil {
		return 0, ce.NewError(ce.ErrCodeCache, ce.ErrMsgInternalServer, tracer, err)
	}

	value, ok := result.(string)
	if !ok {
		return 0, ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, ce.ErrInvalidType)
	}

	authID, err := utils.ToInt64(value)
	if err != nil {
		return 0, ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
	}

	return authID, nil
}
