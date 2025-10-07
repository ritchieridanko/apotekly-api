package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/utils"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errs := ctx.Errors
		if len(errs) == 0 {
			return
		}

		var customErr *ce.Error
		if errors.As(errs[0].Err, &customErr) {

			// TODO (2)

			utils.SetErrorResponse(ctx, customErr.Message, customErr.HTTPStatus())
			return
		}

		// TODO (2)
		utils.SetErrorResponse(ctx, ce.MsgInternalServer, http.StatusInternalServerError)
	}
}
