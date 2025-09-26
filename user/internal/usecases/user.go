package usecases

import (
	"bytes"
	"context"
	"errors"
	"io"
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
	UpdateUser(ctx context.Context, authID int64, data *entities.UserChange) (user *entities.User, err error)
	ChangeProfilePicture(ctx context.Context, authID int64, image multipart.File) (err error)
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

		// upload image if exists
		var pictureURL *string
		userID := utils.GenerateRandomUUID()
		if image != nil {
			imageURL, err := u.uploadImage(ctx, image, userID.String(), "pp", "users", true)
			if err != nil {
				log.Println("WARNING -> failed to upload image:", err.Error())
			} else {
				pictureURL = &imageURL
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
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "GetUser")
	defer span.End()

	return u.ur.GetByAuthID(ctx, authID)
}

func (u *userUsecase) UpdateUser(ctx context.Context, authID int64, data *entities.UserChange) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "UpdateUser")
	defer span.End()

	return u.ur.UpdateUser(ctx, authID, data)
}

func (u *userUsecase) ChangeProfilePicture(ctx context.Context, authID int64, image multipart.File) error {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "ChangeProfilePicture")
	defer span.End()

	userID, err := u.ur.GetUserID(ctx, authID)
	if err != nil {
		return err
	}

	imageURL, err := u.uploadImage(ctx, image, userID.String(), "pp", "users", true)
	if err != nil {
		return err
	}

	return u.ur.UpdateProfilePicture(ctx, authID, imageURL)
}

func (u *userUsecase) uploadImage(ctx context.Context, image multipart.File, publicID, prefix, folder string, overwrite bool) (imageURL string, err error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "uploadImage")
	defer span.End()

	imageBuf, err := io.ReadAll(image)
	if err != nil {
		return "", ce.NewError(span, ce.CodeFileBuffer, ce.MsgInternalServer, err)
	}

	if err := utils.ValidateImageFile(imageBuf); err != nil {
		return "", ce.NewError(span, ce.CodeInvalidPayload, "Invalid file type.", err)
	}

	data := entities.NewUpload{
		File:           bytes.NewReader(imageBuf),
		PublicID:       publicID,
		PublicIDPrefix: prefix,
		Folder:         folder,
		Overwrite:      &overwrite,
	}

	result, err := u.storage.Upload(ctx, &data)
	if err != nil {
		return "", ce.NewError(span, ce.CodeFileUploadFailed, ce.MsgInternalServer, err)
	}

	return result.SecureURL, nil
}
