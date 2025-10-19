package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if strings.TrimSpace(requestID) == "" {
			requestID = utils.NewUUID().String()
		}

		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), constants.CtxKeyRequestID, requestID),
		)

		ctx.Next()
	}
}
