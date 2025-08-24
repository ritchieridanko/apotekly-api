package caches

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
)

func (c *cache) ShouldAccountBeLocked(ctx context.Context, key string) (bool, error) {
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
		ctx,
		c.client,
		[]string{key},
		constants.MaxTotalFailedAuth,
		constants.RedisDurationTotalFailedAuth,
	).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}
