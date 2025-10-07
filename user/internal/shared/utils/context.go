package utils

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/user/internal/shared/constants"
)

func ContextGetAuthID(ctx context.Context) (int64, error) {
	authID, ok := ctx.Value(constants.CtxKeyAuthID).(int64)
	if !ok {
		return 0, errors.New("context has no auth id")
	}
	return authID, nil
}
