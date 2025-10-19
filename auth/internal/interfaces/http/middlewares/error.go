package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/logger"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ErrorHandler(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errs := ctx.Errors
		if len(errs) == 0 {
			return
		}

		var customErr *ce.Error
		if errors.As(errs[0].Err, &customErr) {
			fields := []zapcore.Field{
				zap.String("error_code", string(customErr.Code)),
				zap.String("error_message", customErr.Message),
				zap.String("error_detail", customErr.Error()),
			}

			l.Log(ctx, constants.LogLevelError, "Request Error", customErr.HTTPStatus(), fields...)
			utils.SetErrorResponse(ctx, customErr.Message, customErr.HTTPStatus())
			return
		}

		fields := []zap.Field{
			zap.String("error_detail", errs[0].Err.Error()),
		}

		l.Log(ctx, constants.LogLevelError, "Unhandled Internal Error", http.StatusInternalServerError, fields...)
		utils.SetErrorResponse(ctx, ce.MsgInternalServer, http.StatusInternalServerError)
	}
}
