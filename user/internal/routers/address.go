package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/middlewares"
)

func addressRouters(h handlers.AddressHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.GET("", middlewares.Authenticate(), h.GetAllAddresses)

		rg.POST("", middlewares.Authenticate(), middlewares.RequireVerified(), h.NewAddress)

		rg.DELETE("/:id", middlewares.Authenticate(), middlewares.RequireVerified(), h.DeleteAddress)
	}
}
