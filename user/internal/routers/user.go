package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/middlewares"
)

func userRouters(h handlers.UserHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.GET("/me", middlewares.Authenticate(), h.GetUser)

		rg.POST("", middlewares.Authenticate(), h.NewUser)
	}
}
