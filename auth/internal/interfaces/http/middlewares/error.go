package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/shared/utils"
)

// TODO
// 1: Implement custom logger

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errs := ctx.Errors
		if len(errs) == 0 {
			return
		}

		var customErr *ce.Error
		if errors.As(errs[0].Err, &customErr) {

			// TODO (1)

			utils.SetErrorResponse(ctx, customErr.Message, customErr.HTTPStatus())
			return
		}

		// TODO (1)
		utils.SetErrorResponse(ctx, ce.MsgInternalServer, http.StatusInternalServerError)
	}
}
