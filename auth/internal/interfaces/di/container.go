package di

import (
	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/caches"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/repositories"
	"github.com/ritchieridanko/apotekly-api/auth/internal/app/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infrastructure"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/middlewares"
	"github.com/ritchieridanko/apotekly-api/auth/internal/interfaces/http/router"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/database"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/logger"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/oauth"
)

type Container struct {
	router *router.Router
}

func NewContainer(cfg *configs.Config, infra *infrastructure.Infrastructure) *Container {
	db := database.NewDatabase(infra.DB())
	tx := database.NewTransactor(infra.DB())
	cache := cache.NewCache(infra.Cache(), cfg.Cache.MaxRetries, cfg.Cache.BaseDelay)

	ar := repositories.NewAuthRepository(db)
	oar := repositories.NewOAuthRepository(db)
	sr := repositories.NewSessionRepository(db)

	ac := caches.NewAuthCache(cache)
	oac := caches.NewOAuthCache(cache)
	bcrypt := services.NewBCryptService(cfg.Auth.BCrypt.Cost)
	jwt := services.NewJWTService(&cfg.Auth)
	cookie := services.NewCookieService(cfg.App.Env, true)
	oauth := oauth.Initialize(&cfg.OAuth)
	logger := logger.NewLogger(infra.Logger())

	su := usecases.NewSessionUsecase(sr, tx)
	au := usecases.NewAuthUsecase(ar, ac, su, tx, bcrypt, jwt, cfg)
	oau := usecases.NewOAuthUsecase(oar, ar, oac, ac, su, tx, jwt, cfg)

	ah := handlers.NewAuthHandler(au, cookie, cfg)
	oah := handlers.NewOAuthHandler(oau, au, oauth.Google(), oauth.Microsoft(), cookie, cfg)

	am := middlewares.NewAuthMiddleware(jwt, cfg.App.Name)

	r := router.NewRouter(logger, am, ah, oah, cfg)

	return &Container{router: r}
}

func (c *Container) Router() *router.Router {
	return c.router
}
