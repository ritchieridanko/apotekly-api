package usecases

import (
	"context"
	"errors"

	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"go.opentelemetry.io/otel"
)

// TODO
// (1): Implement image upload to cloud

const userErrorTracer string = "usecase.user"

type UserUsecase interface {
	NewUser(ctx context.Context, authID int64, data *entities.NewUser) (user *entities.User, err error)
}

type userUsecase struct {
	ur repos.UserRepo
	tx db.TxManager
}

func NewUserUsecase(ur repos.UserRepo, tx db.TxManager) UserUsecase {
	return &userUsecase{ur, tx}
}

func (u *userUsecase) NewUser(ctx context.Context, authID int64, data *entities.NewUser) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "NewUser")
	defer span.End()

	var user *entities.User
	err := u.tx.WithTx(ctx, func(ctx context.Context) (err error) {
		exists, err := u.ur.HasUser(ctx, authID)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(span, ce.CodeDBDuplicateData, "User already exists", errors.New("auth id conflict"))
		}

		newData := entities.NewUser{
			UserID:         utils.GenerateRandomUUID(),
			Name:           data.Name,
			Bio:            data.Bio,
			Sex:            data.Sex,
			Birthdate:      data.Birthdate,
			Phone:          data.Phone,
			ProfilePicture: nil, // TODO (1)
		}

		user, err = u.ur.Create(ctx, authID, &newData)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}
