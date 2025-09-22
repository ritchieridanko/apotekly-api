package di

import (
	"database/sql"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/user/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/routers"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/storage"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
)

func SetupDependencies(dbInstance *sql.DB, cloudInstance *cloudinary.Cloudinary) (router *gin.Engine) {
	database := db.NewService(dbInstance)
	txManager := db.NewTxManager(dbInstance)
	storage := storage.NewService(cloudInstance)

	ur := repos.NewUserRepo(database)

	uu := usecases.NewUserUsecase(ur, txManager, storage)

	uh := handlers.NewUserHandler(uu)

	return routers.Initialize(uh)
}
