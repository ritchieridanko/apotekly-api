package utils

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
)

func CtxGetAuthID(ctx context.Context) (int64, error) {
	authID, ok := ctx.Value(constants.CtxKeyAuthID).(int64)
	if !ok {
		return 0, errors.New("auth id not found in request context")
	}
	return authID, nil
}
