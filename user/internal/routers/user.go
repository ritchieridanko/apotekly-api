package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/middlewares"
)

func userRouters(h handlers.UserHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.GET("/me", middlewares.Authenticate(), middlewares.Authorize(), h.GetUser)

		rg.POST("", middlewares.Authenticate(), middlewares.Authorize(), h.NewUser)

		rg.PATCH("/me", middlewares.Authenticate(), middlewares.Authorize(), h.UpdateUser)
		rg.PATCH("/me/profile-picture", middlewares.Authenticate(), middlewares.Authorize(), h.ChangeProfilePicture)
	}
}
