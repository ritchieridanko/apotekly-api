package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
)

func AuthRouters(h handlers.AuthHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.POST("/register", h.Register)
		rg.POST("/login", h.Login)
		rg.POST("/refresh-session", h.RefreshSession)
	}
}
