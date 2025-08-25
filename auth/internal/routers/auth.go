package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/middlewares"
)

func AuthRouters(h handlers.AuthHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.POST("/register", h.Register)
		rg.POST("/login", h.Login)
		rg.POST("/logout", middlewares.Authenticate(), h.Logout)
		rg.POST("/refresh-session", h.RefreshSession)

		rg.PATCH("/email", middlewares.Authenticate(), middlewares.RequireVerified(), h.ChangeEmail)
		rg.PATCH("/password", middlewares.Authenticate(), middlewares.RequireVerified(), h.ChangePassword)
	}
}
