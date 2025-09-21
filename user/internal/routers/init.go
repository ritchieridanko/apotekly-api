package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Initialize(uh handlers.UserHandler) *gin.Engine {
	router := gin.New()

	router.Use(otelgin.Middleware("app.user"))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.ErrorHandler())

	router.ContextWithFallback = true

	api := router.Group("/api/v1")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	user := userRouters(uh)
	user(api.Group("/users"))

	return router
}
