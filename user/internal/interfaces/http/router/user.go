package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/middlewares"
)

type userRoutes struct {
	h    *handlers.UserHandler
	auth *middlewares.AuthMiddleware
}

func newUserRoutes(h *handlers.UserHandler, auth *middlewares.AuthMiddleware) *userRoutes {
	return &userRoutes{h, auth}
}

func (r *userRoutes) register(rg *gin.RouterGroup) {
	rg.GET("/me", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.h.GetUser)
	rg.POST("", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.h.CreateUser)
	rg.PATCH("/me", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.h.UpdateUser)
	rg.PATCH("/me/profile-picture", r.auth.Authenticate(), r.auth.AuthorizeRole(), r.h.ChangeProfilePicture)
}
