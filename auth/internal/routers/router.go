package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/middlewares"
)

func Initialize(auth func(*gin.RouterGroup)) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.CORS())

	router.ContextWithFallback = true

	api := router.Group("/api/v1")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	auth(api.Group("/auth"))

	return router
}
