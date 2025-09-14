package utils

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
)

func GetAuthIDFromContext(ctx context.Context) (authID int64, err error) {
	authID, ok := ctx.Value(constants.RequestKeyAuthID).(int64)
	if !ok {
		return 0, fmt.Errorf("auth id not found in request context")
	}
	return authID, nil
}
