package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/middlewares"
)

func addressRouters(h handlers.AddressHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.GET("", middlewares.Authenticate(), middlewares.Authorize(), h.GetAllAddresses)

		rg.POST("", middlewares.Authenticate(), middlewares.Authorize(), middlewares.RequireVerified(), h.NewAddress)

		rg.PATCH("/:id", middlewares.Authenticate(), middlewares.Authorize(), middlewares.RequireVerified(), h.UpdateAddress)
		rg.PATCH("/:id/primary", middlewares.Authenticate(), middlewares.Authorize(), middlewares.RequireVerified(), h.ChangePrimaryAddress)

		rg.DELETE("/:id", middlewares.Authenticate(), middlewares.Authorize(), middlewares.RequireVerified(), h.DeleteAddress)
	}
}
