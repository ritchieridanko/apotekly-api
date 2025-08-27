package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Initialize(user func(*gin.RouterGroup)) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.ContextWithFallback = true

	api := router.Group("/api/v1")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	user(api.Group("/users"))

	return router
}
