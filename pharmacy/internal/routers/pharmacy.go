package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/middlewares"
)

func pharmacyRouters(h handlers.PharmacyHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.POST("", middlewares.Authenticate(), middlewares.RequireVerified(), h.NewPharmacy)
	}
}
