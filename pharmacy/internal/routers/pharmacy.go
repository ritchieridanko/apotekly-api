package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/middlewares"
)

func pharmacyRouters(h handlers.PharmacyHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.GET("/me", middlewares.Authenticate(), middlewares.RequireVerified(), h.GetPharmacy)

		rg.POST("", middlewares.Authenticate(), middlewares.RequireVerified(), h.NewPharmacy)

		rg.PATCH("/me", middlewares.Authenticate(), middlewares.RequireVerified(), h.UpdatePharmacy)
		rg.PATCH("/me/logo", middlewares.Authenticate(), middlewares.RequireVerified(), h.ChangeLogo)
	}
}
