package di

import (
	"github.com/ritchieridanko/apotekly-api/user/config"
	"github.com/ritchieridanko/apotekly-api/user/internal/infrastructure"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/handlers"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/middlewares"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/router"
	"github.com/ritchieridanko/apotekly-api/user/internal/interfaces/http/validator"
	"github.com/ritchieridanko/apotekly-api/user/internal/repositories"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/database"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/storage"
	"github.com/ritchieridanko/apotekly-api/user/internal/usecases"
)

type Container struct {
	router *router.Router
}

func NewContainer(cfg *config.Config, infra *infrastructure.Infrastructure) *Container {
	db := database.NewDatabase(infra.DB())
	tx := database.NewTransactor(infra.DB())
	storage := storage.NewStorage(infra.Storage())

	ur := repositories.NewUserRepository(db)
	ar := repositories.NewAddressRepository(db)

	uu := usecases.NewUserUsecase(ur, tx, storage)
	au := usecases.NewAddressUsecase(ar, tx)

	v := validator.NewValidator()

	uh := handlers.NewUserHandler(uu, v, int64(cfg.Image.MaxSizeBytes), cfg.Image.AllowedTypes)
	ah := handlers.NewAddressHandler(au, v)

	am := middlewares.NewAuthMiddleware(cfg.App.Name, cfg.JWT.Secret)

	r := router.NewRouter(am, uh, ah, cfg.App.Name)

	return &Container{router: r}
}

func (c *Container) Router() *router.Router {
	return c.router
}
