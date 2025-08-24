package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/auth/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/routers"
	"github.com/ritchieridanko/apotekly-api/auth/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

func SetupDependencies(db *sql.DB) *gin.Engine {
	sr := repos.NewSessionRepo(db)
	ar := repos.NewAuthRepo(db)

	txManager := dbtx.New(db)

	su := usecases.NewSessionUsecase(sr)
	au := usecases.NewAuthUsecase(ar, txManager, su)

	ah := handlers.NewAuthHandler(au)

	return routers.Initialize(routers.AuthRouters(ah))
}
