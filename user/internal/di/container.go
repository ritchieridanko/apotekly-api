package di

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/routers"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
)

func SetupDependencies(dbInstance *sql.DB) (router *gin.Engine) {
	database := db.NewService(dbInstance)
	txManager := db.NewTxManager(dbInstance)

	ur := repos.NewUserRepo(database)

	uu := usecases.NewUserUsecase(ur, txManager)

	uh := handlers.NewUserHandler(uu)

	return routers.Initialize(uh)
}
