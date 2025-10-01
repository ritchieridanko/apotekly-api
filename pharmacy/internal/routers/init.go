package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Initialize(ph handlers.PharmacyHandler) *gin.Engine {
	router := gin.New()

	router.Use(otelgin.Middleware("app.pharmacy"))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.ErrorHandler())

	router.ContextWithFallback = true

	api := router.Group("/api/v1")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	pharmacy := pharmacyRouters(ph)
	pharmacy(api.Group("/pharmacies"))

	return router
}
