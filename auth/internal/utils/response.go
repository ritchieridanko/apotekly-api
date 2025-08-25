package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/dto"
)

func SetResponse(ctx *gin.Context, message string, data any, code int) {
	response := dto.Response{
		Message: message,
		Data:    data,
	}
	ctx.JSON(code, &response)
}

func SetErrorResponse(ctx *gin.Context, message string, code int) {
	response := dto.Response{
		Message: message,
	}
	ctx.AbortWithStatusJSON(code, &response)
}
