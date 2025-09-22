package usecases

import (
	"bytes"
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/ritchieridanko/apotekly-api/user/internal/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/repos"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/db"
	"github.com/ritchieridanko/apotekly-api/user/internal/services/storage"
	"github.com/ritchieridanko/apotekly-api/user/internal/utils"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "usecase.user"

type UserUsecase interface {
	NewUser(ctx context.Context, authID int64, data *entities.NewUser, image multipart.File) (user *entities.User, err error)
	GetUser(ctx context.Context, authID int64) (user *entities.User, err error)
}

type userUsecase struct {
	ur      repos.UserRepo
	tx      db.TxManager
	storage storage.StorageService
}

func NewUserUsecase(ur repos.UserRepo, tx db.TxManager, storage storage.StorageService) UserUsecase {
	return &userUsecase{ur, tx, storage}
}

func (u *userUsecase) NewUser(ctx context.Context, authID int64, data *entities.NewUser, image multipart.File) (*entities.User, error) {
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

		var pictureURL *string
		userID := utils.GenerateRandomUUID()
		if image != nil {
			imageBuf, err := utils.FileProcessImage(image)
			if err != nil {
				log.Println("FAIL -> failed to process image:", err)
			} else {
				overwrite := true
				uploadData := entities.NewUpload{
					File:           bytes.NewReader(imageBuf),
					PublicID:       userID.String(),
					PublicIDPrefix: "pp",
					Folder:         "users",
					Overwrite:      &overwrite,
				}

				result, err := u.storage.Upload(ctx, &uploadData)
				if err != nil {
					log.Println("FAIL -> failed to upload image:", err)
				} else {
					pictureURL = &result.SecureURL
				}
			}
		}

		newData := entities.NewUser{
			UserID:         userID,
			Name:           data.Name,
			Bio:            data.Bio,
			Sex:            data.Sex,
			Birthdate:      data.Birthdate,
			Phone:          data.Phone,
			ProfilePicture: pictureURL,
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

func (u *userUsecase) GetUser(ctx context.Context, authID int64) (*entities.User, error) {
	return u.ur.GetByAuthID(ctx, authID)
}
