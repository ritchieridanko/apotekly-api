package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
)

const cacheErrorTracer string = "service.cache"

type CacheService interface {
	Has(ctx context.Context, key string) (exists bool, err error)
	Del(ctx context.Context, keys ...string) (err error)

	ShouldLockAccount(ctx context.Context, authID int64) (shouldBe bool, err error)
	NewPasswordResetToken(ctx context.Context, authID int64, token string, duration time.Duration) (err error)
	ConsumePasswordResetToken(ctx context.Context, token string) (authID int64, err error)
	NewVerificationToken(ctx context.Context, authID int64, token string, duration time.Duration) (err error)
	ConsumeVerificationToken(ctx context.Context, token string) (authID int64, err error)
	NewEmailChangeToken(ctx context.Context, authID int64, newEmail, token string, duration time.Duration) (err error)
	ConsumeEmailChangeToken(ctx context.Context, token string) (authID int64, newEmail string, err error)
}

type cacheService struct {
	client     *redis.Client
	maxRetries int
	baseDelay  int
}

func NewService(client *redis.Client) CacheService {
	return &cacheService{client, config.CacheGetMaxRetries(), config.CacheGetBaseDelay()}
}

func (cs *cacheService) Has(ctx context.Context, key string) (bool, error) {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "Has")
	defer span.End()

	var lastErr error
	for attempt := 0; attempt <= cs.maxRetries; attempt++ {
		result, err := cs.client.Exists(ctx, key).Result()
		if err == nil {
			return result > 0, nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}

		if err := backoffWait(ctx, cs.baseDelay, attempt); err != nil {
			return false, ce.NewError(span, ce.CodeCacheBackoffWait, ce.MsgInternalServer, err)
		}
	}

	return false, ce.NewError(span, ce.CodeCacheQueryExecution, ce.MsgInternalServer, lastErr)
}

func (cs *cacheService) Del(ctx context.Context, keys ...string) error {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "Del")
	defer span.End()

	var lastErr error
	for attempt := 0; attempt <= cs.maxRetries; attempt++ {
		err := cs.client.Del(ctx, keys...).Err()
		if err == nil {
			return nil
		}

		lastErr = err
		if !isRetryable(err) {
			break
		}

		if err := backoffWait(ctx, cs.baseDelay, attempt); err != nil {
			return ce.NewError(span, ce.CodeCacheBackoffWait, ce.MsgInternalServer, err)
		}
	}

	return ce.NewError(span, ce.CodeCacheQueryExecution, ce.MsgInternalServer, lastErr)
}

func (cs *cacheService) ShouldLockAccount(ctx context.Context, authID int64) (bool, error) {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "ShouldLockAccount")
	defer span.End()

	key := utils.CacheCreateKey(constants.CacheKeyTotalFailedAuth, authID)

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

	result, err := script.Run(
		ctx, cs.client, []string{key},
		constants.CacheMaxTotalFailedAuth,
		constants.CacheDurationTotalFailedAuth*60,
	).Int()
	if err != nil {
		return false, ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	return result == 1, nil
}

func (cs *cacheService) NewPasswordResetToken(ctx context.Context, authID int64, token string, duration time.Duration) error {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "NewPasswordResetToken")
	defer span.End()

	authKey := utils.CacheCreateKey(constants.CacheKeyPasswordResetAuth, authID)
	tokenKey := utils.CacheCreateKey(constants.CacheKeyPasswordResetToken, token)

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

	_, err := script.Run(
		ctx, cs.client,
		[]string{authKey, tokenKey, constants.CacheKeyPasswordResetToken},
		token, authID, int(duration.Seconds()),
	).Result()
	if err != nil {
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (cs *cacheService) ConsumePasswordResetToken(ctx context.Context, token string) (int64, error) {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "ConsumePasswordResetToken")
	defer span.End()

	tokenKey := utils.CacheCreateKey(constants.CacheKeyPasswordResetToken, token)

	script := redis.NewScript(`
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`)

	result, err := script.Run(ctx, cs.client, []string{tokenKey, constants.CacheKeyPasswordResetAuth}).Result()
	if err == redis.Nil {
		return 0, ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, err)
	}
	if err != nil {
		return 0, ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	value, ok := result.(string)
	if !ok {
		return 0, ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, ce.ErrTypeAssertionFailed)
	}

	authID, err := utils.ToInt64(value)
	if err != nil {
		return 0, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, ce.ErrTypeConversionFailed)
	}

	return authID, nil
}

func (cs *cacheService) NewVerificationToken(ctx context.Context, authID int64, token string, duration time.Duration) error {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "NewVerificationToken")
	defer span.End()

	authKey := utils.CacheCreateKey(constants.CacheKeyVerificationAuth, authID)
	tokenKey := utils.CacheCreateKey(constants.CacheKeyVerificationToken, token)

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

	_, err := script.Run(
		ctx, cs.client,
		[]string{authKey, tokenKey, constants.CacheKeyVerificationToken},
		token, authID, int(duration.Seconds()),
	).Result()
	if err != nil {
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (cs *cacheService) ConsumeVerificationToken(ctx context.Context, token string) (int64, error) {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "ConsumeVerificationToken")
	defer span.End()

	tokenKey := utils.CacheCreateKey(constants.CacheKeyVerificationToken, token)

	script := redis.NewScript(`
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`)

	result, err := script.Run(ctx, cs.client, []string{tokenKey, constants.CacheKeyVerificationAuth}).Result()
	if err == redis.Nil {
		return 0, ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, err)
	}
	if err != nil {
		return 0, ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	value, ok := result.(string)
	if !ok {
		return 0, ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, ce.ErrTypeAssertionFailed)
	}

	authID, err := utils.ToInt64(value)
	if err != nil {
		return 0, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, ce.ErrTypeConversionFailed)
	}

	return authID, nil
}

func (cs *cacheService) NewEmailChangeToken(ctx context.Context, authID int64, newEmail, token string, duration time.Duration) error {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "NewEmailChangeToken")
	defer span.End()

	authKey := utils.CacheCreateKey(constants.CacheKeyEmailChangeAuth, authID)
	tokenKey := utils.CacheCreateKey(constants.CacheKeyEmailChangeToken, token)

	script := redis.NewScript(`
		local oldToken = redis.call("GET", KEYS[1])
		if oldToken then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. oldToken)
		end
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[4])
		redis.call("HSET", KEYS[2], "authID", ARGV[2], "newEmail", ARGV[3])
		redis.call("EXPIRE", KEYS[2], ARGV[4])
		return 1
	`)

	_, err := script.Run(
		ctx, cs.client,
		[]string{authKey, tokenKey, constants.CacheKeyEmailChangeToken},
		token, authID, newEmail, int(duration.Seconds()),
	).Result()
	if err != nil {
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	return nil
}

func (cs *cacheService) ConsumeEmailChangeToken(ctx context.Context, token string) (int64, string, error) {
	ctx, span := otel.Tracer(cacheErrorTracer).Start(ctx, "ConsumeEmailChangeToken")
	defer span.End()

	tokenKey := utils.CacheCreateKey(constants.CacheKeyEmailChangeToken, token)

	script := redis.NewScript(`
		local data = redis.call("HMGET", KEYS[1], "authID", "newEmail")
		if data and data[1] then
			local authID = data[1]
			local newEmail = data[2]
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return {authID, newEmail}
		end
		return nil
	`)

	result, err := script.Run(ctx, cs.client, []string{tokenKey, constants.CacheKeyEmailChangeAuth}).Result()
	if err == redis.Nil || result == nil {
		return 0, "", ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, err)
	}
	if err != nil {
		return 0, "", ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, err)
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return 0, "", ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, ce.ErrTypeAssertionFailed)
	}

	authStr, ok1 := values[0].(string)
	newEmail, ok2 := values[1].(string)
	if !ok1 || !ok2 {
		return 0, "", ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, ce.ErrTypeAssertionFailed)
	}

	authID, err := utils.ToInt64(authStr)
	if err != nil {
		return 0, "", ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, ce.ErrTypeConversionFailed)
	}

	return authID, newEmail, nil
}
