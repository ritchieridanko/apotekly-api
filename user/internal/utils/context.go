package utils

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/user/internal/constants"
)

func ContextGetAuthID(ctx context.Context) (authID int64, err error) {
	authID, ok := ctx.Value(constants.CtxKeyAuthID).(int64)
	if !ok {
		return 0, errors.New("auth id not found in request context")
	}
	return authID, nil
}
