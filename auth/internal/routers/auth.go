package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
)

func AuthRouters(h handlers.AuthHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.POST("", h.Register)
	}
}
