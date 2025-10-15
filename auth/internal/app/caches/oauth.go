package caches

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"go.opentelemetry.io/otel"
)

const oAuthErrorTracer string = "cache.oauth"

type OAuthCache interface {
	StoreAuth(ctx context.Context, code string, auth *entities.Auth, duration time.Duration) (err error)
	GetAuth(ctx context.Context, code string) (auth *entities.Auth, err error)
}

type oAuthCache struct {
	cache *cache.Cache
}

func NewOAuthCache(cache *cache.Cache) OAuthCache {
	return &oAuthCache{cache}
}

func (c *oAuthCache) StoreAuth(ctx context.Context, code string, auth *entities.Auth, duration time.Duration) error {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "StoreAuth")
	defer span.End()

	codeKey := fmt.Sprintf("%s:%s", constants.CachePrefixOAuthStore, code)
	isVerified := strconv.FormatBool(auth.IsVerified)
	createdAt := auth.CreatedAt.Format(time.RFC3339)
	updatedAt := auth.UpdatedAt.Format(time.RFC3339)

	// id: authID, em: email, rid: roleID, iv: isVerified, cat: createdAt, uat: updatedAt
	script := `
		redis.call("DEL", KEYS[1])
		redis.call("HSET", KEYS[1],
			"id", ARGV[1],
			"em", ARGV[2],
			"rid", ARGV[3],
			"iv", ARGV[4],
			"cat", ARGV[5],
			"uat", ARGV[6]
		)
		redis.call("EXPIRE", KEYS[1], ARGV[7])
		return 1
	`

	_, err := c.cache.Evaluate(
		ctx, "hs:oasa", script, []string{codeKey},
		auth.ID, auth.Email, auth.RoleID, isVerified,
		createdAt, updatedAt, int(duration.Seconds()),
	)
	if err != nil {
		wErr := fmt.Errorf("failed to store auth: %w", err)
		return ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	return nil
}

func (c *oAuthCache) GetAuth(ctx context.Context, code string) (*entities.Auth, error) {
	ctx, span := otel.Tracer(oAuthErrorTracer).Start(ctx, "GetAuth")
	defer span.End()

	codeKey := fmt.Sprintf("%s:%s", constants.CachePrefixOAuthStore, code)

	// id: authID, em: email, rid: roleID, iv: isVerified, cat: createdAt, uat: updatedAt
	script := `
		local data = redis.call("HMGET", KEYS[1], "id", "em", "rid", "iv", "cat", "uat")
		if data and data[1] then
			local authID = data[1]
			local email = data[2]
			local roleID = data[3]
			local isVerified = data[4]
			local createdAt = data[5]
			local updatedAt = data[6]
			redis.call("DEL", KEYS[1])
			return {authID, email, roleID, isVerified, createdAt, updatedAt}
		end
		return nil
	`

	result, err := c.cache.Evaluate(ctx, "hs:oaga", script, []string{codeKey})
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		if errors.Is(err, ce.ErrCacheNil) {
			return nil, ce.NewError(span, ce.CodeCacheValueNotFound, ce.MsgInvalidToken, wErr)
		}
		return nil, ce.NewError(span, ce.CodeCacheScriptExecution, ce.MsgInternalServer, wErr)
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 6 {
		err := fmt.Errorf("failed to fetch auth: %w", ce.ErrTypeAssertionFailed)
		return nil, ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, err)
	}

	authStr, ok1 := values[0].(string)
	email, ok2 := values[1].(string)
	roleStr, ok3 := values[2].(string)
	isVerifiedStr, ok4 := values[3].(string)
	createdAtStr, ok5 := values[4].(string)
	updatedAtStr, ok6 := values[5].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		err := fmt.Errorf("failed to fetch auth: %w", ce.ErrTypeAssertionFailed)
		return nil, ce.NewError(span, ce.CodeTypeAssertionFailed, ce.MsgInternalServer, err)
	}

	authID, err := utils.ToInt64(authStr)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return nil, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	roleID, err := utils.ToInt64(roleStr)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return nil, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	isVerified, err := strconv.ParseBool(isVerifiedStr)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return nil, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return nil, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		wErr := fmt.Errorf("failed to fetch auth: %w", err)
		return nil, ce.NewError(span, ce.CodeTypeConversionFailed, ce.MsgInternalServer, wErr)
	}

	auth := entities.Auth{
		ID:         authID,
		Email:      email,
		RoleID:     int16(roleID),
		IsVerified: isVerified,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	return &auth, nil
}
