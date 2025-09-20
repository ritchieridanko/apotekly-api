package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Initialize(ah handlers.AuthHandler, oah handlers.OAuthHandler) *gin.Engine {
	router := gin.New()

	router.Use(otelgin.Middleware("app.auth"))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.CORS())

	router.ContextWithFallback = true

	api := router.Group("/api/v1")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	auth := authRouters(ah)
	auth(api.Group("/auth"))

	oauth := oAuthRouters(oah)
	oauth(api.Group("/oauth"))

	return router
}
