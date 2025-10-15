package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/handlers"
)

type oAuthRouter struct {
	h *handlers.OAuthHandler
}

func newOAuthRouter(h *handlers.OAuthHandler) *oAuthRouter {
	return &oAuthRouter{h}
}

func (r *oAuthRouter) register(rg *gin.RouterGroup) {
	rg.POST("/exchange", r.h.OAuthExchange)

	google := rg.Group("/google")
	google.GET("", r.h.GoogleOAuth)
	google.GET("/callback", r.h.GoogleCallback)

	microsoft := rg.Group("/microsoft")
	microsoft.GET("", r.h.MicrosoftOAuth)
	microsoft.GET("/callback", r.h.MicrosoftCallback)
}
