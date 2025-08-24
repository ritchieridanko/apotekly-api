package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly-api/auth/internal/caches"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/routers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

func SetupDependencies(db *sql.DB, redis *redis.Client) (router *gin.Engine) {
	sr := repos.NewSessionRepo(db)
	ar := repos.NewAuthRepo(db)

	txManager := dbtx.NewTxManager(db)
	cache := caches.NewCache(redis)

	su := usecases.NewSessionUsecase(sr, txManager)
	au := usecases.NewAuthUsecase(ar, txManager, su, cache)

	ah := handlers.NewAuthHandler(au)

	return routers.Initialize(routers.AuthRouters(ah))
}
