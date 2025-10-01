package di

import (
	"database/sql"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/handlers"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/repos"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/routers"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/services/storage"
	"github.com/ritchieridanko/apotekly-api/pharmacy/internal/usecases"
)

func SetupDependencies(dbInstance *sql.DB, cloudInstance *cloudinary.Cloudinary) (router *gin.Engine) {
	database := db.NewService(dbInstance)
	txManager := db.NewTxManager(dbInstance)
	storage := storage.NewService(cloudInstance)

	pr := repos.NewPharmacyRepo(database)

	pu := usecases.NewPharmacyUsecase(pr, txManager, storage)

	ph := handlers.NewPharmacyHandler(pu)

	return routers.Initialize(ph)
}
