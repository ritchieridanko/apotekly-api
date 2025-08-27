package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"github.com/ritchieridanko/apotekly-api/user/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/user/pkg/dbtx"
)

// TODO
// (1): Implement Image Upload to Cloud

const UserErrorTracer = ce.UserUsecaseTracer

type UserUsecase interface {
	NewUser(ctx context.Context, authID int64, data *entities.NewUser) (err error)
}

type userUsecase struct {
	ur repos.UserRepo
	tx dbtx.TxManager
}

func NewUserUsecase(ur repos.UserRepo, tx dbtx.TxManager) UserUsecase {
	return &userUsecase{ur, tx}
}

func (u *userUsecase) NewUser(ctx context.Context, authID int64, data *entities.NewUser) error {
	tracer := UserErrorTracer + ": NewUser()"

	return u.tx.ReturnError(ctx, func(ctx context.Context) (err error) {
		exists, err := u.ur.HasUser(ctx, authID)
		if err != nil {
			return err
		}
		if exists {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgUserAlreadyExists, tracer, ce.ErrUserAlreadyExists)
		}

		if data.Sex != nil && !utils.ValidatorIsSexValid(*data.Sex) {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidDataSex, tracer, ce.ErrInvalidDataSex)
		}
		if data.Birthdate != nil && !utils.ValidatorIsBirthdateValid(*data.Birthdate) {
			return ce.NewError(ce.ErrCodeInvalidAction, ce.ErrMsgInvalidDataBirthdate, tracer, ce.ErrInvalidDataBirthdate)
		}

		newData := entities.NewUser{
			UserID:         utils.GenerateRandomUUID(),
			Name:           data.Name,
			Bio:            data.Bio,
			Sex:            data.Sex,
			Birthdate:      data.Birthdate,
			Phone:          data.Phone,
			ProfilePicture: data.ProfilePicture, // TODO (1)
		}

		return u.ur.Create(ctx, authID, &newData)
	})
}
