package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dtos"
)

func SetResponse(ctx *gin.Context, message string, data any, code int) {
	response := dtos.Response{
		Message: message,
		Data:    data,
	}
	ctx.JSON(code, &response)
}
