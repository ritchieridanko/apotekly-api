package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
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
			log.Printf("ERROR -> %s: %s", customErr.Code, customErr.Err)
			utils.SetErrorResponse(ctx, customErr.Message, customErr.ToExternalErrorCode())
			return
		}

		log.Println("ERROR ->", errs[0].Err)
		utils.SetErrorResponse(ctx, ce.MsgInternalServer, http.StatusInternalServerError)
	}
}
