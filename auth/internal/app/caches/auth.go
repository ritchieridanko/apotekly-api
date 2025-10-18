package caches

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/services/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"go.opentelemetry.io/otel"
)

const authErrorTracer string = "cache.auth"

type AuthCache interface {
	CreateResetToken(ctx context.Context, authID int64, token string, duration time.Duration) (err error)
	UseResetToken(ctx context.Context, token string) (authID int64, err error)
	CreateVerificationToken(ctx context.Context, authID int64, token string, duration time.Duration) (err error)
	UseVerificationToken(ctx context.Context, token string) (authID int64, err error)
	CreateEmailChangeToken(ctx context.Context, authID int64, newEmail, token string, duration time.Duration) (err error)
	UseEmailChangeToken(ctx context.Context, token string) (authID int64, newEmail string, err error)
	UnreserveEmail(ctx context.Context, email string) (err error)
	ResetTokenExists(ctx context.Context, token string) (exists bool, err error)
	IsEmailReserved(ctx context.Context, email string) (exists bool, err error)
}

type authCache struct {
	cache *cache.Cache
}

func NewAuthCache(cache *cache.Cache) AuthCache {
	return &authCache{cache}
}

func (c *authCache) CreateResetToken(ctx context.Context, authID int64, token string, duration time.Duration) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "CreateResetToken")
	defer span.End()

	authKey := fmt.Sprintf("%s:%d", constants.CachePrefixReset, authID)
	tokenKey := fmt.Sprintf("%s:%s", constants.CachePrefixReset, token)

	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`

	_, err := c.cache.Evaluate(
		ctx, "hs:crt", script,
		[]string{authKey, tokenKey, constants.CachePrefixReset},
		token, strconv.FormatInt(authID, 10), int(duration.Seconds()),
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create reset token: %w", err)
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	return nil
}

func (c *authCache) UseResetToken(ctx context.Context, token string) (int64, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UseResetToken")
	defer span.End()

	tokenKey := fmt.Sprintf("%s:%s", constants.CachePrefixReset, token)

	script := `
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`

	result, err := c.cache.Evaluate(ctx, "hs:urt", script, []string{tokenKey, constants.CachePrefixReset})
	if err != nil {
		wErr := fmt.Errorf("failed to use reset token: %w", err)
		if errors.Is(err, ce.ErrCacheNil) {
			return 0, ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, wErr)
		}
		return 0, ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	authID, err := utils.ToInt64Any(result)
	if err != nil {
		wErr := fmt.Errorf("failed to use reset token: %w", err)
		return 0, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	return authID, nil
}

func (c *authCache) CreateVerificationToken(ctx context.Context, authID int64, token string, duration time.Duration) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "CreateVerificationToken")
	defer span.End()

	authKey := fmt.Sprintf("%s:%d", constants.CachePrefixVerification, authID)
	tokenKey := fmt.Sprintf("%s:%s", constants.CachePrefixVerification, token)

	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[3])
		redis.call("SET", KEYS[2], ARGV[2], "EX", ARGV[3])
		return 1
	`

	_, err := c.cache.Evaluate(
		ctx, "hs:cvt", script,
		[]string{authKey, tokenKey, constants.CachePrefixVerification},
		token, strconv.FormatInt(authID, 10), int(duration.Seconds()),
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create verification token: %w", err)
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	return nil
}

func (c *authCache) UseVerificationToken(ctx context.Context, token string) (int64, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UseVerificationToken")
	defer span.End()

	tokenKey := fmt.Sprintf("%s:%s", constants.CachePrefixVerification, token)

	script := `
		local authID = redis.call("GET", KEYS[1])
		if authID then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return authID
		end
		return nil
	`

	result, err := c.cache.Evaluate(ctx, "hs:uvt", script, []string{tokenKey, constants.CachePrefixVerification})
	if err != nil {
		wErr := fmt.Errorf("failed to use verification token: %w", err)
		if errors.Is(err, ce.ErrCacheNil) {
			return 0, ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, wErr)
		}
		return 0, ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	authID, err := utils.ToInt64Any(result)
	if err != nil {
		wErr := fmt.Errorf("failed to use verification token: %w", err)
		return 0, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	return authID, nil
}

func (c *authCache) CreateEmailChangeToken(ctx context.Context, authID int64, newEmail, token string, duration time.Duration) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "CreateEmailChangeToken")
	defer span.End()

	authKey := fmt.Sprintf("%s:%d", constants.CachePrefixEmailChange, authID)
	tokenKey := fmt.Sprintf("%s:%s", constants.CachePrefixEmailChange, token)
	emailKey := fmt.Sprintf("%s:%s", constants.CachePrefixEmailReservation, newEmail)

	// id: authID, ne: newEmail
	script := `
		local token = redis.call("GET", KEYS[1])
		if token then
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[3] .. ":" .. token)
		end
		local reserved = redis.call("SET", KEYS[4], ARGV[2], "NX", "EX", ARGV[4])
		if reserved == nil then
			return 0
		end
		redis.call("EXPIRE", KEYS[4], ARGV[4])
		redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[4])
		redis.call("HSET", KEYS[2], "id", ARGV[2], "ne", ARGV[3])
		redis.call("EXPIRE", KEYS[2], ARGV[4])
		return 1
	`

	res, err := c.cache.Evaluate(
		ctx, "hs:cect", script,
		[]string{authKey, tokenKey, constants.CachePrefixEmailChange, emailKey},
		token, strconv.FormatInt(authID, 10), newEmail, int(duration.Seconds()),
	)
	if err != nil {
		wErr := fmt.Errorf("failed to create email change token: %w", err)
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}
	if res.(int64) == 0 {
		err := fmt.Errorf("failed to create email change token: %w", errors.New("email already reserved"))
		return ce.NewError(span, ce.CodeAuthEmailConflict, "Email is already registered", err)
	}

	return nil
}

func (c *authCache) UseEmailChangeToken(ctx context.Context, token string) (int64, string, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UseEmailChangeToken")
	defer span.End()

	tokenKey := fmt.Sprintf("%s:%s", constants.CachePrefixEmailChange, token)

	// id: authID, ne: newEmail
	script := `
		local data = redis.call("HMGET", KEYS[1], "id", "ne")
		if data and data[1] then
			local authID = data[1]
			local email = data[2]
			redis.call("DEL", KEYS[1])
			redis.call("DEL", KEYS[2] .. ":" .. authID)
			return {authID, email}
		end
		return nil
	`

	result, err := c.cache.Evaluate(ctx, "hs:uect", script, []string{tokenKey, constants.CachePrefixEmailChange})
	if err != nil {
		wErr := fmt.Errorf("failed to use email change token: %w", err)
		if errors.Is(err, ce.ErrCacheNil) {
			return 0, "", ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, wErr)
		}
		return 0, "", ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		err := fmt.Errorf("failed to use email change token: %w", ce.ErrTypeAssertionFailed)
		return 0, "", ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, err)
	}

	authStr, ok1 := values[0].(string)
	newEmail, ok2 := values[1].(string)
	if !ok1 || !ok2 {
		err := fmt.Errorf("failed to use email change token: %w", ce.ErrTypeAssertionFailed)
		return 0, "", ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, err)
	}

	authID, err := utils.ToInt64(authStr)
	if err != nil {
		wErr := fmt.Errorf("failed to use email change token: %w", err)
		return 0, "", ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	return authID, newEmail, nil
}

func (c *authCache) UnreserveEmail(ctx context.Context, email string) error {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "UnreserveEmail")
	defer span.End()

	key := fmt.Sprintf("%s:%s", constants.CachePrefixEmailReservation, email)

	if err := c.cache.Delete(ctx, key); err != nil {
		wErr := fmt.Errorf("failed to unreserve email: %w", err)
		return ce.NewError(span, ce.CodeCacheQueryExecution, ce.MsgInternalServer, wErr)
	}

	return nil
}

func (c *authCache) ResetTokenExists(ctx context.Context, token string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "ResetTokenExists")
	defer span.End()

	key := fmt.Sprintf("%s:%s", constants.CachePrefixReset, token)

	exists, err := c.cache.Exists(ctx, key)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch reset token: %w", err)
		return false, ce.NewError(span, ce.CodeCacheQueryExecution, ce.MsgInternalServer, wErr)
	}

	return exists, nil
}

func (c *authCache) IsEmailReserved(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.Tracer(authErrorTracer).Start(ctx, "IsEmailReserved")
	defer span.End()

	key := fmt.Sprintf("%s:%s", constants.CachePrefixEmailReservation, email)

	exists, err := c.cache.Exists(ctx, key)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch email: %w", err)
		return false, ce.NewError(span, ce.CodeCacheQueryExecution, ce.MsgInternalServer, wErr)
	}

	return exists, nil
}
