package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/mailer"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/routers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/cache"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/email"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/oauth"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
)

func SetupDependencies(dbInstance *sql.DB, redisClient *redis.Client, mailer mailer.Mailer) (router *gin.Engine) {
	database := db.NewService(dbInstance)
	txManager := db.NewTxManager(dbInstance)
	cache := cache.NewService(redisClient)
	email := email.NewService(mailer)

	sr := repos.NewSessionRepo(database)
	ar := repos.NewAuthRepo(database)
	oar := repos.NewOAuthRepo(database)

	su := usecases.NewSessionUsecase(sr, txManager)
	au := usecases.NewAuthUsecase(ar, su, txManager, cache, email)
	oau := usecases.NewOAuthUsecase(oar, ar, su, txManager, cache, email)

	google, microsoft := oauth.Initialize()

	ah := handlers.NewAuthHandler(au)
	oah := handlers.NewOAuthHandler(oau, au, google, microsoft)

	return routers.Initialize(ah, oah)
}
