package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/middlewares"
)

type addressRoutes struct {
	h    *handlers.AddressHandler
	auth *middlewares.AuthMiddleware
}

func newAddressRoutes(h *handlers.AddressHandler, auth *middlewares.AuthMiddleware) *addressRoutes {
	return &addressRoutes{h, auth}
}

func (r *addressRoutes) register(rg *gin.RouterGroup) {
	rg.GET("", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.h.GetAllAddresses)
	rg.POST("", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.auth.AuthorizeVerification(), r.h.CreateAddress)
	rg.PATCH("/:id", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.auth.AuthorizeVerification(), r.h.UpdateAddress)
	rg.PATCH("/:id/primary", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.auth.AuthorizeVerification(), r.h.SetPrimaryAddress)
	rg.DELETE("/:id", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.auth.AuthorizeVerification(), r.h.DeleteAddress)
}
