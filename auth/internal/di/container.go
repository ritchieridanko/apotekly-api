package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/internal/caches"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/mailer"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/routers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/services/email"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

func SetupDependencies(db *sql.DB, redis *redis.Client, mailer *mailer.Mailer) (router *gin.Engine) {
	sr := repos.NewSessionRepo(db)
	ar := repos.NewAuthRepo(db)
	oar := repos.NewOAuthRepo(db)

	txManager := dbtx.NewTxManager(db)
	cache := caches.NewCache(redis)
	email := email.NewEmailService(mailer)

	su := usecases.NewSessionUsecase(sr, txManager)
	oau := usecases.NewOAuthUsecase(oar)
	au := usecases.NewAuthUsecase(ar, txManager, su, oau, cache, email)

	ah := handlers.NewAuthHandler(au)

	return routers.Initialize(routers.AuthRouters(ah))
}
