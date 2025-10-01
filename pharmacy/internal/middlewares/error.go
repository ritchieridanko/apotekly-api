package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/ce"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/utils"
)

// TODO
// 1: When custom logger is made, log the error here

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
			utils.SetErrorResponse(ctx, customErr.Message, customErr.ToExternalErrorCode())
			return
		}

		// TODO (1)
		utils.SetErrorResponse(ctx, ce.MsgInternalServer, http.StatusInternalServerError)
	}
}
