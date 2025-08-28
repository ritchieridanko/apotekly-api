package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/routers"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
	"github.com/ritchieridanko/apotekly-api/user/pkg/dbtx"
)

func SetupDependencies(db *sql.DB) (router *gin.Engine) {
	ur := repos.NewUserRepo(db)

	txManager := dbtx.NewTxManager(db)

	uu := usecases.NewUserUsecase(ur, txManager)

	uh := handlers.NewUserHandler(uu)

	return routers.Initialize(routers.UserRouters(uh))
}
