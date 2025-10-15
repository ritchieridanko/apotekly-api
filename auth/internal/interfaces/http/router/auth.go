package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/middlewares"
)

type authRouter struct {
	h    *handlers.AuthHandler
	auth *middlewares.AuthMiddleware
}

func newAuthRouter(h *handlers.AuthHandler, auth *middlewares.AuthMiddleware) *authRouter {
	return &authRouter{h, auth}
}

func (r *authRouter) register(rg *gin.RouterGroup) {
	rg.GET("/email/available", r.h.IsEmailRegistered)
	rg.GET("/verify-account/confirm", r.h.VerifyAccount)
	rg.GET("/change-email/confirm", r.h.ConfirmEmailChange)

	rg.POST("/register", r.h.Register)
	rg.POST("/login", r.h.Login)
	rg.POST("/logout", r.auth.Authenticate(), r.h.Logout)
	rg.POST("/refresh-session", r.h.RefreshSession)
	rg.POST("/forgot-password", r.h.ForgotPassword)
	rg.POST("/reset-password/confirm", r.h.ResetPassword)
	rg.POST("/reset-password/validate", r.h.IsResetTokenValid)
	rg.POST("/verify-account/resend", r.auth.Authenticate(), r.h.ResendVerification)

	rg.PATCH("/change-email/request", r.auth.Authenticate(), r.auth.RequireVerified(), r.h.ChangeEmail)
	rg.PATCH("/password", r.auth.Authenticate(), r.auth.RequireVerified(), r.h.ChangePassword)
}
