package usecases

import (
	"context"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/constants"
	"github.com/ritchieridanko/apotekly-api/auth/internal/entities"
	"github.com/ritchieridanko/apotekly-api/auth/internal/repos"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/dbtx"
)

const AuthErrorTracer = ce.AuthUsecaseTracer

type AuthUsecase interface {
	Register(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (token *entities.AuthToken, err error)
}

type authUsecase struct {
	ar repos.AuthRepo
	tx dbtx.TxManager
	su SessionUsecase
}

func NewAuthUsecase(ar repos.AuthRepo, tx dbtx.TxManager, su SessionUsecase) AuthUsecase {
	return &authUsecase{ar, tx, su}
}

func (u *authUsecase) Register(ctx context.Context, data *entities.NewAuth, request *entities.NewRequest) (*entities.AuthToken, error) {
	tracer := AuthErrorTracer + ": Register()"
	currentTime := time.Now().UTC()
	sessionDuration := time.Duration(config.GetSessionDuration()) * time.Minute

	result, err := u.tx.ReturnAnyAndError(ctx, func(ctx context.Context) (any, error) {
		normalizedEmail := utils.NormalizeString(data.Email)

		exists, err := u.ar.IsEmailRegistered(ctx, normalizedEmail)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ce.NewError(ce.ErrCodeConflict, ce.ErrMsgEmailRegistered, tracer, ce.ErrEmailAlreadyRegistered)
		}

		hashedPassword, err := utils.HashPassword(data.Password)
		if err != nil {
			return nil, ce.NewError(ce.ErrCodePassHashing, ce.ErrMsgInternalServer, tracer, err)
		}

		newData := entities.NewAuth{
			Email:    normalizedEmail,
			Password: hashedPassword,
			Role:     constants.RoleCustomer,
		}

		authID, err := u.ar.Create(ctx, &newData)
		if err != nil {
			return nil, err
		}

		sessionToken := utils.GenerateRandomToken()
		accessToken, err := utils.GenerateJWTToken(authID, data.Role, false)
		if err != nil {
			return nil, ce.NewError(ce.ErrCodeToken, ce.ErrMsgInternalServer, tracer, err)
		}

		sessionData := entities.NewSession{
			AuthID:    authID,
			Token:     sessionToken,
			UserAgent: request.UserAgent,
			IPAddress: request.IPAddress,
			ExpiresAt: currentTime.Add(sessionDuration),
		}

		if err := u.su.NewSessionOnRegister(ctx, &sessionData); err != nil {
			return nil, err
		}

		token := entities.AuthToken{
			AccessToken:  accessToken,
			SessionToken: sessionToken,
		}

		return token, nil
	})
	if err != nil {
		return nil, err
	}

	token, ok := result.(entities.AuthToken)
	if !ok {
		return nil, ce.NewError(ce.ErrCodeInvalidType, ce.ErrMsgInternalServer, tracer, ce.ErrInvalidType)
	}

	return &token, nil
}
