package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/logger"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"go.uber.org/zap"
)

func Logger(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now().UTC()
		ctx.Next()

		latency := zap.String("latency", time.Since(start).String())
		statusCode := ctx.Writer.Status()

		if statusCode < http.StatusBadRequest {
			l.Log(ctx, constants.LogLevelInfo, "Request", statusCode, latency)
		} else {
			l.Log(ctx, constants.LogLevelWarn, "Request Warning", statusCode, latency)
		}
	}
}
