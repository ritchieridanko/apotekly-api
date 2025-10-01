package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/handlers"
)

// TODO
// 1: Add Authenticate() and RequireVerified() middlewares

func pharmacyRouters(h handlers.PharmacyHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		rg.POST("", h.NewPharmacy) // TODO (1)
	}
}
