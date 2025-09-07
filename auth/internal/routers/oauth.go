package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
)

func OAuthRouters(h handlers.OAuthHandler) func(*gin.RouterGroup) {
	return func(rg *gin.RouterGroup) {
		google := rg.Group("/google")
		google.GET("", h.GoogleOAuth)
		google.GET("/callback", h.GoogleCallback)
	}
}
