package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func NewRouter(
	am *middlewares.AuthMiddleware,
	ah *handlers.AuthHandler,
	oah *handlers.OAuthHandler,
	cfg *configs.Config,
) *Router {
	r := gin.New()

	r.Use(otelgin.Middleware(cfg.App.Name))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.CORS(cfg.Client.BaseURL))

	r.ContextWithFallback = true

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	auth := newAuthRouter(ah, am)
	auth.register(api.Group("/auth"))

	oauth := newOAuthRouter(oah)
	oauth.register(api.Group("/oauth"))

	return &Router{router: r}
}

func (r *Router) Engine() *gin.Engine {
	return r.router
}
