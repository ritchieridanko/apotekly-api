package middlewares

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"github.com/ritchieridanko/apotekly-api/user/pkg/ce"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors[0]

		var cErr *ce.Error
		if errors.As(err.Err, &cErr) {
			log.Printf("ERROR: (%d) => %v", cErr.Code, cErr.Err)
			utils.SetErrorResponse(ctx, cErr.Message, ce.MapToExternalErrorCode(cErr.Code))
			return
		}

		log.Println("ERROR:", err.Err)
		utils.SetErrorResponse(ctx, ce.ErrMsgInternalServer, http.StatusInternalServerError)
	}
}
