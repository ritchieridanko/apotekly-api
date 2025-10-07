package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/ritchieridanko/apotekly-api/user/internal/entities"
	"github.com/ritchieridanko/apotekly-api/user/internal/repositories"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/database"
	"github.com/ritchieridanko/apotekly-api/user/internal/service/storage"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/ce"
	"github.com/ritchieridanko/apotekly-api/user/internal/shared/utils"
	"go.opentelemetry.io/otel"
)

const userErrorTracer string = "usecase.user"

type UserUsecase interface {
	CreateUser(ctx context.Context, authID int64, data *entities.CreateUser, image multipart.File) (createdUser *entities.User, err error)
	GetUser(ctx context.Context, authID int64) (user *entities.User, err error)
	UpdateUser(ctx context.Context, authID int64, data *entities.UpdateUser) (updatedUser *entities.User, err error)
	ChangeProfilePicture(ctx context.Context, authID int64, image multipart.File) (user *entities.User, err error)
}

type userUsecase struct {
	ur         repositories.UserRepository
	transactor *database.Transactor
	storage    *storage.Storage
}

func NewUserUsecase(
	ur repositories.UserRepository,
	transactor *database.Transactor,
	storage *storage.Storage,
) UserUsecase {
	return &userUsecase{ur, transactor, storage}
}

func (u *userUsecase) CreateUser(ctx context.Context, authID int64, data *entities.CreateUser, image multipart.File) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "CreateUser")
	defer span.End()

	var user *entities.User
	err := u.transactor.WithTx(ctx, func(ctx context.Context) (err error) {
		exists, err := u.ur.Exists(ctx, authID)
		if err != nil {
			return err
		}
		if exists {
			err := fmt.Errorf("failed to create user: %w", errors.New("conflict auth id"))
			return ce.NewError(span, ce.CodeDBDuplicateData, "User already exists", err)
		}

		// upload image if exists
		var profilePicture *string
		userID := utils.NewUUID()
		if image != nil {
			imageURL, err := u.uploadImage(ctx, image, userID.String(), "pp", "users/profile_pictures", true)
			if err != nil {
				log.Println("WARNING -> ", err.Error())
			} else {
				profilePicture = &imageURL
			}
		}

		newData := entities.CreateUser{
			ID:             userID,
			Name:           data.Name,
			Bio:            data.Bio,
			Sex:            data.Sex,
			Birthdate:      data.Birthdate,
			Phone:          data.Phone,
			ProfilePicture: profilePicture,
		}

		user, err = u.ur.Create(ctx, authID, &newData)
		return err
	})

	return user, err
}

func (u *userUsecase) GetUser(ctx context.Context, authID int64) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "GetUser")
	defer span.End()

	return u.ur.GetByAuthID(ctx, authID)
}

func (u *userUsecase) UpdateUser(ctx context.Context, authID int64, data *entities.UpdateUser) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "UpdateUser")
	defer span.End()

	return u.ur.Update(ctx, authID, data)
}

func (u *userUsecase) ChangeProfilePicture(ctx context.Context, authID int64, image multipart.File) (*entities.User, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "ChangeProfilePicture")
	defer span.End()

	userID, err := u.ur.GetUserID(ctx, authID)
	if err != nil {
		return nil, err
	}

	profilePicture, err := u.uploadImage(ctx, image, userID.String(), "pp", "users/profile_pictures", true)
	if err != nil {
		return nil, err
	}

	return u.ur.UpdateProfilePicture(ctx, authID, profilePicture)
}

func (u *userUsecase) uploadImage(ctx context.Context, image multipart.File, publicID, prefix, folder string, overwrite bool) (string, error) {
	ctx, span := otel.Tracer(userErrorTracer).Start(ctx, "uploadImage")
	defer span.End()

	data := entities.UploadParams{
		File:      image,
		PublicID:  publicID,
		Prefix:    prefix,
		Folder:    folder,
		Overwrite: &overwrite,
	}

	result, err := u.storage.Upload(ctx, &data)
	if err != nil {
		wErr := fmt.Errorf("failed to upload image: %w", err)
		return "", ce.NewError(span, ce.CodeFileUploadFailed, ce.MsgInternalServer, wErr)
	}

	return result.SecureURL, nil
}
