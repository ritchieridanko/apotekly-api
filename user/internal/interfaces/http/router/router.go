package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func NewRouter(
	am *middlewares.AuthMiddleware,
	uh *handlers.UserHandler,
	ah *handlers.AddressHandler,

	appName string,
) *Router {
	r := gin.New()

	r.Use(otelgin.Middleware(appName))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.ErrorHandler())

	r.ContextWithFallback = true

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	user := newUserRoutes(uh, am)
	user.register(api.Group("/users"))

	address := newAddressRoutes(ah, am)
	address.register(api.Group("/addresses"))

	return &Router{router: r}
}

func (r *Router) Engine() *gin.Engine {
	return r.router
}
